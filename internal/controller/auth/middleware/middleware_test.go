package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"timeline/internal/config"
	"timeline/internal/entity"
	"timeline/internal/libs/jwtlib"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type MiddlewareTestSuite struct {
	suite.Suite
	Middeware      *Middleware
	tokenCfg       config.Token
	mockPrivateKey *rsa.PrivateKey
}

func (suite *MiddlewareTestSuite) SetupTest() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.mockPrivateKey = privateKey
	suite.tokenCfg = config.Token{AccessTTL: 10 * time.Minute, RefreshTTL: 10 * time.Minute}
	suite.Middeware = New(
		suite.mockPrivateKey,
		zap.NewExample(),
	)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suiter := &MiddlewareTestSuite{}
	suiter.SetupTest()
	suite.Run(t, suiter)
}

func (suite *MiddlewareTestSuite) TestValidToken() {
	metadata := &entity.TokenMetadata{ID: 0, IsOrg: false}
	validToken, err := jwtlib.NewToken(suite.mockPrivateKey, suite.tokenCfg, metadata, "access")
	suite.Require().NoError(err)

	r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewBufferString(""))
	r.Header.Set("Authorization", "Bearer "+validToken)
	w := httptest.NewRecorder()

	token, err := suite.Middeware.ExtractToken(w, r)
	suite.NoError(err)
	suite.NotNil(token)
}

func (suite *MiddlewareTestSuite) TestEmptyAuthHeader() {
	r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewBufferString(""))
	r.Header.Set("Authorization", "")
	w := httptest.NewRecorder()

	_, err := suite.Middeware.ExtractToken(w, r)
	suite.Equal(ErrAuthHeaderEmpty, err)
}

func (suite *MiddlewareTestSuite) TestTokenNotFound() {
	r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewBufferString(""))
	r.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()

	_, err := suite.Middeware.ExtractToken(w, r)
	suite.Equal(ErrTokenNotFound, err)
}

func (suite *MiddlewareTestSuite) TestSuspectToken() {
	metadata := &entity.TokenMetadata{ID: 0, IsOrg: false}
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	suite.Require().NoError(err)
	suspectToken, err := jwtlib.NewToken(privateKey, suite.tokenCfg, metadata, "access")
	suite.NoError(err)

	r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewBufferString(""))
	r.Header.Set("Authorization", "Bearer "+suspectToken)
	w := httptest.NewRecorder()

	_, err = suite.Middeware.ExtractToken(w, r)
	suite.Error(err)
}
