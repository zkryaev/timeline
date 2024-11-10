package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"timeline/internal/config"
	jwtlib "timeline/internal/libs/jwt"
	"timeline/internal/model"
	"timeline/internal/model/dto"
	"timeline/internal/usecase/auth/validation"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Repository interface {
	SaveUser(ctx context.Context, user *model.User) (int, error)
	User(ctx context.Context) (*model.User, error)
	SaveOrg(ctx context.Context, org *model.Organization) (int, error)
	Organization(ctx context.Context) (*model.Organization, error)
}

// TODO: в Usecase надо логировать

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthUseCase struct {
	secret   *rsa.PrivateKey
	Repo     Repository
	TokenCfg config.Token
	Logger   *zap.Logger
}

func New(key *rsa.PrivateKey, storage Repository, cfg config.Token, logger *zap.Logger) *AuthUseCase {
	return &AuthUseCase{
		secret:   key,
		Repo:     storage,
		TokenCfg: cfg,
		Logger:   logger,
	}
}

func (a *AuthUseCase) Login(ctx context.Context, req dto.LoginReq) (*dto.TokenPair, error) {
	// идем в базу и проверяем существует ли пользователь с таким email
	_, err := a.Repo.User(ctx)
	if err != nil {
		return nil, err
	}

	// если ДА, то пароль полученного пользователя хеш входного пароля сравнивается с тем что вернула БД
	// если да, возвращаем токен
	return nil, nil
}

func (a *AuthUseCase) UserRegister(ctx context.Context, req dto.UserRegisterReq) (*dto.TokenPair, error) {
	// идем в базу и проверяем существует ли пользователь с таким email
	// если НЕТ, то добавляем его в бд
	// возвращаем токен
	return nil, nil
}

func (a *AuthUseCase) OrgRegister(ctx context.Context, req dto.OrgRegisterReq) (*dto.TokenPair, error) {
	// идем в базу и проверяем существует ли пользователь с таким email
	// если НЕТ, то добавляем его в бд
	// возвращаем токен
	return nil, nil
}

func (a *AuthUseCase) UpdateRefreshToken(ctx context.Context, req *jwt.Token) (string, error) {
	// Валидируем Claims токена. Есть ли они и нормальные ли.
	err := validation.ValidateTokenClaims(req)
	if err != nil {
		return "", err
	}
	// Здесь уже спокойно кастую если выше проблем не возникло
	metadata := &model.TokenMetadata{
		ID:    req.Claims.(jwt.MapClaims)["id"].(uint64),
		IsOrg: req.Claims.(jwt.MapClaims)["is_org"].(bool),
	}
	token, err := jwtlib.NewToken(a.secret, a.TokenCfg, metadata, req.Claims.(jwt.MapClaims)["type"].(string))
	if err != nil {
		return "", err
	}
	return token, nil
}
