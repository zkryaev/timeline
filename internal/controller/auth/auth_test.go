package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"timeline/internal/config"
	"timeline/internal/controller/scope"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/sugar/jwtlib"
	mocks "timeline/mocks/controller"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type AuthTestSuite struct {
	suite.Suite
	Auth            *AuthCtrl
	mockAuthUseCase *mocks.AuthUseCase
	mockMiddleware  *mocks.Middleware
	mockPrivateKey  *rsa.PrivateKey
	tokenCfg        config.Token
	appcfg          config.Application
}

func (suite *AuthTestSuite) SetupTest() {
	suite.mockAuthUseCase = &mocks.AuthUseCase{}
	suite.mockMiddleware = &mocks.Middleware{}
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	suite.mockPrivateKey = privateKey
	suite.tokenCfg = config.Token{AccessTTL: 10 * time.Minute, RefreshTTL: 10 * time.Minute}
	suite.Auth = New(suite.mockAuthUseCase, suite.mockMiddleware, zap.NewExample(), scope.NewDefaultSettings(suite.appcfg))
}

func TestAuthTestSuite(t *testing.T) {
	suiter := &AuthTestSuite{}
	suiter.SetupTest()
	suite.Run(t, suiter)
}

func (suite *AuthTestSuite) TestLoginSuccess() {
	correctJSON := `{
		"email": "test@test.ru",
		"password": "testpasswd1241",
		"is_org": false
	}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(correctJSON))

	pair, err := jwtlib.NewTokenPair(suite.mockPrivateKey, suite.tokenCfg, &entity.TokenData{ID: 0, IsOrg: false})
	suite.Require().NoError(err)

	input := authdto.LoginReq{
		Credentials: authdto.Credentials{Email: "test@test.ru", Password: "testpasswd1241"},
		IsOrg:       false,
	}
	logger := suite.Auth.Logger.With(zap.String("uuid", ""))
	suite.mockAuthUseCase.On("Login", r.Context(), logger, &input).Return(pair, nil)
	suite.Auth.Login(w, r)
	suite.Equal(http.StatusOK, w.Result().StatusCode, w.Body.String())
}

func (suite *AuthTestSuite) TestUpdateAccessTokenRefreshError() {
	access, err := jwtlib.NewToken(suite.mockPrivateKey, suite.tokenCfg, &entity.TokenData{ID: 0, IsOrg: false}, "access")
	suite.Require().NoError(err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(access))
	r.Header.Set("Authorization", "Bearer "+access)

	token, err := jwt.Parse(access, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Join(errors.New(""), fmt.Errorf("token sign: %s", token.Method.Alg()))
		}
		publicKey := &suite.mockPrivateKey.PublicKey
		return publicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))
	suite.Require().NoError(err)

	suite.mockMiddleware.On("ExtractToken", r).Return(token, nil)
	logger := suite.Auth.Logger.With(zap.String("uuid", ""))
	suite.mockAuthUseCase.On("UpdateAccessToken", r.Context(), logger, token).Return(nil, errors.New("mustn't be called"))
	suite.Auth.PutAccessToken(w, r)
	suite.Require().Equal(http.StatusBadRequest, w.Result().StatusCode)
}
