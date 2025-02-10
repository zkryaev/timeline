package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"
	"timeline/internal/config"
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/infrastructure/models"
	"timeline/internal/libs/passwd"
	mocks "timeline/mocks/infrastructure"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type AuthUseCaseTestSuite struct {
	suite.Suite
	authUseCase  *AuthUseCase
	mockUserRepo *mocks.UserRepository
	mockOrgRepo  *mocks.OrgRepository
	mockCodeRepo *mocks.CodeRepository
	mockMailRepo *mocks.Mail
}

func (suite *AuthUseCaseTestSuite) SetupTest() {
	suite.mockUserRepo = &mocks.UserRepository{}
	suite.mockOrgRepo = &mocks.OrgRepository{}
	suite.mockCodeRepo = &mocks.CodeRepository{}
	suite.mockMailRepo = &mocks.Mail{}
	mockPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		suite.T().Fatal(err.Error())
	}
	suite.authUseCase = New(
		mockPrivateKey,
		suite.mockUserRepo,
		suite.mockOrgRepo,
		suite.mockCodeRepo,
		suite.mockMailRepo,
		config.Token{AccessTTL: 10 * time.Minute, RefreshTTL: 10 * time.Minute},
		zap.NewExample(),
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
	require.NoError(suite.T(), err)

	expectedExpiration := &models.ExpInfo{
		ID:        1,
		Verified:  true,
		CreatedAt: time.Now(),
		Hash:      hash,
	}

	suite.mockCodeRepo.On("AccountExpiration", mock.Anything, email, isOrg).Return(expectedExpiration, nil)

	// Вызов метода Login
	req := &authdto.LoginReq{
		Credentials: authdto.Credentials{
			Email:    email,
			Password: password,
		},
		IsOrg: isOrg,
	}
	tokens, err := suite.authUseCase.Login(ctx, req)
	suite.NoError(err)
	suite.NotNil(tokens)
	suite.NotNil(tokens.AccessToken)
	suite.NotNil(tokens.RefreshToken)

	// Проверка вызова моков
	suite.mockCodeRepo.AssertExpectations(suite.T())
}

// func (suite *AuthUseCaseTestSuite) TestUserRegister(t *testing.T) {

// }

// func (suite *AuthUseCaseTestSuite) TestOrgRegister(t *testing.T) {

// }

// func (suite *AuthUseCaseTestSuite) TestVerifyCode(t *testing.T) {

// }
