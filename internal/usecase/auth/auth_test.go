package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"
	"timeline/internal/config"
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/infrastructure/mapper/codemap"
	"timeline/internal/infrastructure/models"
	"timeline/internal/libs/passwd"
	mocks "timeline/mocks/infrastructure"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type AuthUseCaseTestSuite struct {
	suite.Suite
	logger         *zap.Logger
	authUseCase    *AuthUseCase
	mockUserRepo   *mocks.UserRepository
	mockOrgRepo    *mocks.OrgRepository
	mockCodeRepo   *mocks.CodeRepository
	mockMailRepo   *mocks.Mail
	mockPrivateKey *rsa.PrivateKey
}

func (suite *AuthUseCaseTestSuite) SetupTest() {
	suite.mockUserRepo = &mocks.UserRepository{}
	suite.mockOrgRepo = &mocks.OrgRepository{}
	suite.mockCodeRepo = &mocks.CodeRepository{}
	suite.mockMailRepo = &mocks.Mail{}
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.mockPrivateKey = privateKey
	suite.logger = zap.NewExample()
	suite.authUseCase = New(
		suite.mockPrivateKey,
		suite.mockUserRepo,
		suite.mockOrgRepo,
		suite.mockCodeRepo,
		suite.mockMailRepo,
		config.Token{AccessTTL: 10 * time.Minute, RefreshTTL: 10 * time.Minute},
	)
}

func TestAuthUseCaseTestSuite(t *testing.T) {
	suiter := &AuthUseCaseTestSuite{}
	suiter.SetupTest()
	suite.Run(t, suiter)
}

func (suite *AuthUseCaseTestSuite) TestLoginSuccess() {
	ctx := context.Background()
	email := "test@example.com"
	password := "SecurePassword123"
	isOrg := false

	hash, err := passwd.GetHash(password)
	suite.Require().NoError(err)

	expectedExpiration := &models.ExpInfo{
		ID:        0,
		Verified:  true,
		CreatedAt: time.Now(),
		Hash:      hash,
	}

	suite.mockCodeRepo.On("AccountExpiration", mock.Anything, email, isOrg).Return(expectedExpiration, nil)

	req := &authdto.LoginReq{
		Credentials: authdto.Credentials{
			Email:    email,
			Password: password,
		},
		IsOrg: isOrg,
	}
	tokens, err := suite.authUseCase.Login(ctx, suite.logger, req)
	suite.NoError(err)
	suite.NotNil(tokens)
	suite.NotNil(tokens.AccessToken)
	suite.NotNil(tokens.RefreshToken)

	suite.mockCodeRepo.AssertExpectations(suite.T())
}

func (suite *AuthUseCaseTestSuite) TestLoginFail() {
	ctx := context.Background()

	req := &authdto.LoginReq{
		Credentials: authdto.Credentials{
			Email:    "test@example.com",
			Password: "SecurePassword123",
		},
		IsOrg: false,
	}

	hash, err := passwd.GetHash("another_passwd")
	suite.Require().NoError(err)

	expectedExpiration := &models.ExpInfo{
		ID:        0,
		Verified:  true,
		CreatedAt: time.Now(),
		Hash:      hash,
	}

	suite.mockCodeRepo.On("AccountExpiration", mock.Anything, req.Email, req.IsOrg).Return(expectedExpiration, nil)

	tokens, err := suite.authUseCase.Login(ctx, suite.logger, req)
	suite.Error(err)
	suite.Nil(tokens)

	suite.mockCodeRepo.AssertExpectations(suite.T())
}

func (suite *AuthUseCaseTestSuite) TestLoginNotVerifiedAccountExpired() {
	ctx := context.Background()

	req := &authdto.LoginReq{
		Credentials: authdto.Credentials{
			Email:    "test@example.com",
			Password: "SecurePassword123",
		},
		IsOrg: false,
	}

	hash, err := passwd.GetHash(req.Password)
	suite.Require().NoError(err)

	expectedExpiration := &models.ExpInfo{
		ID:        0,
		Verified:  false,
		CreatedAt: time.Time{},
		Hash:      hash,
	}

	suite.mockCodeRepo.On("AccountExpiration", mock.Anything, req.Email, req.IsOrg).Return(expectedExpiration, nil)

	tokens, err := suite.authUseCase.Login(ctx, suite.logger, req)
	suite.Equal(ErrAccountExpired, err)
	suite.Nil(tokens)

	suite.mockCodeRepo.AssertExpectations(suite.T())
}

func (suite *AuthUseCaseTestSuite) TestLoginVerifiedAccountExpired() {
	ctx := context.Background()
	req := &authdto.LoginReq{
		Credentials: authdto.Credentials{
			Email:    "test@example.com",
			Password: "SecurePassword123",
		},
		IsOrg: false,
	}

	hash, err := passwd.GetHash(req.Password)
	suite.Require().NoError(err)

	expectedExpiration := &models.ExpInfo{
		ID:        0,
		Verified:  true,
		CreatedAt: time.Time{},
		Hash:      hash,
	}

	suite.mockCodeRepo.On("AccountExpiration", mock.Anything, req.Email, req.IsOrg).Return(expectedExpiration, nil)

	tokens, err := suite.authUseCase.Login(ctx, suite.logger, req)
	suite.NoError(err)
	suite.NotNil(tokens)

	suite.mockCodeRepo.AssertExpectations(suite.T())
}

func (suite *AuthUseCaseTestSuite) TestVerifyCodeFresh() {
	ctx := context.Background()
	req := &authdto.VerifyCodeReq{
		Code: "0000",
	}

	notExpired := time.Now().Add(5 * time.Minute)

	suite.mockCodeRepo.On("VerifyCode", mock.Anything, codemap.ToModel(req)).Return(notExpired, nil)
	suite.mockCodeRepo.On("ActivateAccount", mock.Anything, req.ID, req.IsOrg).Return(nil)

	tokens, err := suite.authUseCase.VerifyCode(ctx, suite.logger, req)
	suite.NoError(err)
	suite.NotNil(tokens)

	suite.mockCodeRepo.AssertExpectations(suite.T())
}

func (suite *AuthUseCaseTestSuite) TestVerifyCodeExpired() {
	ctx := context.Background()
	req := &authdto.VerifyCodeReq{
		Code: "0000",
	}

	notExpired := time.Now().Add(-10 * time.Minute)

	suite.mockCodeRepo.On("VerifyCode", mock.Anything, codemap.ToModel(req)).Return(notExpired, nil)

	tokens, err := suite.authUseCase.VerifyCode(ctx, suite.logger, req)
	suite.Equal(ErrCodeExpired, err)
	suite.Nil(tokens)

	suite.mockCodeRepo.AssertExpectations(suite.T())
}

func (suite *AuthUseCaseTestSuite) TestUpdateAccessTokenSuccess() {
	ctx := context.Background()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":     1.0,
		"is_org": false,
		"type":   "access",
		"exp":    time.Now(),
	})

	tokens, err := suite.authUseCase.UpdateAccessToken(ctx, suite.logger, token)
	suite.NoError(err)
	suite.NotNil(tokens)
}

func (suite *AuthUseCaseTestSuite) TestUpdateAccessTokenBadClaims() {
	ctx := context.Background()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":     1.0,
		"is_org": false,
		"type":   00000000000,
		"exp":    time.Now(),
	})

	tokens, err := suite.authUseCase.UpdateAccessToken(ctx, suite.logger, token)
	suite.Error(err)
	suite.Nil(tokens)
}
