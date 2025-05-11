package middleware

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"timeline/internal/config"
	"timeline/internal/controller/monitoring/metrics"
	"timeline/internal/controller/scope"
	"timeline/internal/entity"
	"timeline/internal/sugar/jwtlib"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MiddlewareTestSuite struct {
	suite.Suite
	Middeware      Middleware
	tokenCfg       config.Token
	appCfg         config.Application
	mockPrivateKey *rsa.PrivateKey
}

func (suite *MiddlewareTestSuite) SetupTest() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	suite.mockPrivateKey = privateKey
	suite.tokenCfg = config.Token{AccessTTL: 10 * time.Minute, RefreshTTL: 10 * time.Minute}
	suite.appCfg.Settings.EnableAuthorization = true
	suite.appCfg.Settings.EnableMedia = false
	suite.appCfg.Settings.EnableMail = false
	suite.appCfg.Settings.EnableMetrics = false
	var metricList *metrics.Metrics
	settings := scope.NewDefaultSettings(suite.appCfg)
	suite.Middeware = New(
		suite.mockPrivateKey,
		zap.New(zapcore.NewNopCore()),
		scope.NewDefaultRoutes(settings),
		metricList,
	)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suiter := &MiddlewareTestSuite{}
	suiter.SetupTest()
	suite.Run(t, suiter)
}

func (suite *MiddlewareTestSuite) TestValidToken() {
	metadata := &entity.TokenData{ID: 0, IsOrg: false}
	validToken, err := jwtlib.NewToken(suite.mockPrivateKey, suite.tokenCfg, metadata, "access")
	suite.Require().NoError(err)

	r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewBufferString(""))
	r.Header.Set("Authorization", "Bearer "+validToken)

	token, err := suite.Middeware.ExtractToken(r)
	suite.NoError(err)
	suite.Require().NotNil(token)
}

func (suite *MiddlewareTestSuite) TestEmptyAuthHeader() {
	r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewBufferString(""))
	r.Header.Set("Authorization", "")

	_, err := suite.Middeware.ExtractToken(r)
	suite.Equal(ErrAuthHeaderEmpty, err)
}

func (suite *MiddlewareTestSuite) TestTokenNotFound() {
	r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewBufferString(""))
	r.Header.Set("Authorization", "Bearer ")

	_, err := suite.Middeware.ExtractToken(r)
	suite.Equal(ErrTokenNotFound, err)
}

func (suite *MiddlewareTestSuite) TestSuspectToken() {
	metadata := &entity.TokenData{ID: 0, IsOrg: false}
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	suite.Require().NoError(err)
	suspectToken, err := jwtlib.NewToken(privateKey, suite.tokenCfg, metadata, "access")
	suite.NoError(err)

	r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewBufferString(""))
	r.Header.Set("Authorization", "Bearer "+suspectToken)

	_, err = suite.Middeware.ExtractToken(r)
	suite.Error(err)
}

func (suite *MiddlewareTestSuite) TestAnotherTokenSign() {
	metadata := &entity.TokenData{ID: 0, IsOrg: false}
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	suite.Require().NoError(err)
	suspectToken, err := jwtlib.NewToken(privateKey, suite.tokenCfg, metadata, "access")
	suite.NoError(err)
	suite.NotEmpty(len(suspectToken), 0)

	r := httptest.NewRequest(http.MethodGet, "/test", bytes.NewBufferString(""))
	r.Header.Set("Authorization", "Bearer "+suspectToken)

	_, err = suite.Middeware.ExtractToken(r)
	suite.Contains(err.Error(), jwt.ErrTokenSignatureInvalid.Error())
}
