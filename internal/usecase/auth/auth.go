package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"timeline/internal/config"
	"timeline/internal/entity"
	"timeline/internal/entity/dto"
	jwtlib "timeline/internal/libs/jwtlib"
	"timeline/internal/libs/passwd"
	"timeline/internal/repository/models"
	"timeline/internal/usecase/auth/validation"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type UserRepository interface {
	UserSave(ctx context.Context, user *entity.UserInfo, creds *entity.Credentials) (int, error)
	UserByEmail(ctx context.Context, email string) (*entity.User, *entity.Credentials, error)
	UserByID(ctx context.Context, id int) (*entity.User, *entity.Credentials, error)
	UserIsExist(ctx context.Context, email string) (*models.IsExistResponse, error)
	UserSaveCode(ctx context.Context, code string, user_id int) error
	UserCode(ctx context.Context, code string, user_id int) (string, error)
}
type OrgRepository interface {
	OrgSave(ctx context.Context, org *entity.OrgInfo, creds *entity.Credentials, cityName string) (int, error)
	OrgByEmail(ctx context.Context, email string) (*entity.Organization, *entity.Credentials, error)
	OrgByID(ctx context.Context, id int) (*entity.Organization, *entity.Credentials, error)
	OrgIsExist(ctx context.Context, email string) (*models.IsExistResponse, error)
	OrgSaveCode(ctx context.Context, code string, org_id int) error
	OrgCode(ctx context.Context, code string, org_id int) (string, error)
}

// TODO: в Usecase надо логировать

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountExpired     = errors.New("account expired")
)

type AuthUseCase struct {
	secret   *rsa.PrivateKey
	user     UserRepository
	org      OrgRepository
	TokenCfg config.Token
	Logger   *zap.Logger
}

func New(key *rsa.PrivateKey, userRepo UserRepository, orgRepo OrgRepository, cfg config.Token, logger *zap.Logger) *AuthUseCase {
	return &AuthUseCase{
		secret:   key,
		user:     userRepo,
		org:      orgRepo,
		TokenCfg: cfg,
		Logger:   logger,
	}
}

// Идем в БД и проверяем существует ли пользователь с таким email и не сгорел ли его аккаунт,
// если ДА получаем его, сравниваем хеш паролей если совпадают генерим и отдаем токены
func (a *AuthUseCase) Login(ctx context.Context, req dto.LoginReq) (*dto.TokenPair, error) {
	var MetaInfo *models.IsExistResponse
	var err error
	switch req.IsOrg {
	case false: // user
		MetaInfo, err = a.user.UserIsExist(ctx, req.Email)
	case true: // org
		MetaInfo, err = a.org.OrgIsExist(ctx, req.Email)
	}
	if validation.IsAccountExpired(MetaInfo.CreatedAt) {
		return nil, ErrAccountExpired
	}
	if err != nil {
		a.Logger.Error(
			"failed login account",
			zap.Error(err),
		)
	}
	if err := passwd.CompareWithHash(req.Password, MetaInfo.Hash); err != nil {
		a.Logger.Error(
			"failed login account",
			zap.Error(err),
		)
		return nil, err
	}

	tokens, err := jwtlib.NewTokenPair(
		a.secret,
		a.TokenCfg,
		&entity.TokenMetadata{
			ID:    uint64(MetaInfo.ID),
			IsOrg: req.IsOrg,
		},
	)
	if err != nil {
		a.Logger.Error(
			"failed login account",
			zap.Error(err),
		)
		return nil, err
	}
	return tokens, nil
}

// идем в БД и проверяем существует ли пользователь с таким email,
// если НЕТ, то добавляем его в БД, если ДА -> ошибка
func (a *AuthUseCase) UserRegister(ctx context.Context, req dto.UserRegisterReq) (*dto.TokenPair, error) {

	return nil, nil
}

// Аналогично работе с пользователем
func (a *AuthUseCase) OrgRegister(ctx context.Context, req dto.OrgRegisterReq) (*dto.TokenPair, error) {

	return nil, nil
}

func (a *AuthUseCase) UpdateRefreshToken(ctx context.Context, req *jwt.Token) (string, error) {
	// Валидируем Claims токена. Есть ли они и нормальные ли.
	err := validation.ValidateTokenClaims(req)
	if err != nil {
		return "", err
	}
	// Здесь уже спокойно кастую если выше проблем не возникло
	metadata := &entity.TokenMetadata{
		ID:    req.Claims.(jwt.MapClaims)["id"].(uint64),
		IsOrg: req.Claims.(jwt.MapClaims)["is_org"].(bool),
	}
	token, err := jwtlib.NewToken(a.secret, a.TokenCfg, metadata, req.Claims.(jwt.MapClaims)["type"].(string))
	if err != nil {
		return "", err
	}
	return token, nil
}
