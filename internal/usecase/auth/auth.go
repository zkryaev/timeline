package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"time"
	"timeline/internal/config"
	"timeline/internal/entity"
	"timeline/internal/entity/dto"
	jwtlib "timeline/internal/libs/jwtlib"
	"timeline/internal/libs/passwd"
	"timeline/internal/libs/verification"
	"timeline/internal/repository/mail/notify"
	"timeline/internal/repository/models"
	"timeline/internal/usecase/auth/validation"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type UserRepository interface {
	UserSave(ctx context.Context, user *entity.UserInfo, creds *entity.Credentials) (int, error)
	UserByEmail(ctx context.Context, email string) (*entity.User, *entity.Credentials, error)
	UserByID(ctx context.Context, id int) (*entity.User, *entity.Credentials, error)
	UserGetMetaInfo(ctx context.Context, email string) (*models.MetaInfo, error)
	UserSaveCode(ctx context.Context, code string, user_id int) error
	UserCode(ctx context.Context, code string, user_id int) (time.Time, error)
	UserActivateAccount(ctx context.Context, user_id int) error
	UserIsExist(ctx context.Context, email string) (int, error)
}
type OrgRepository interface {
	OrgSave(ctx context.Context, org *entity.OrgInfo, creds *entity.Credentials, cityName string) (int, error)
	OrgByEmail(ctx context.Context, email string) (*entity.Organization, *entity.Credentials, error)
	OrgByID(ctx context.Context, id int) (*entity.Organization, *entity.Credentials, error)
	OrgGetMetaInfo(ctx context.Context, email string) (*models.MetaInfo, error)
	OrgSaveCode(ctx context.Context, code string, org_id int) error
	OrgCode(ctx context.Context, code string, org_id int) (time.Time, error)
	OrgActivateAccount(ctx context.Context, user_id int) error
	OrgIsExist(ctx context.Context, email string) (int, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountExpired     = errors.New("account expired")
	ErrCodeExpired        = errors.New("code expired")
)

type AuthUseCase struct {
	secret   *rsa.PrivateKey
	user     UserRepository
	org      OrgRepository
	mail     notify.Mail
	TokenCfg config.Token
	Logger   *zap.Logger
}

func New(key *rsa.PrivateKey, userRepo UserRepository, orgRepo OrgRepository, mailSrv notify.Mail, cfg config.Token, logger *zap.Logger) *AuthUseCase {
	return &AuthUseCase{
		secret:   key,
		user:     userRepo,
		org:      orgRepo,
		mail:     mailSrv,
		TokenCfg: cfg,
		Logger:   logger,
	}
}

// Идем в БД и проверяем существует ли пользователь с таким email и не сгорел ли его аккаунт,
// если ДА получаем его, сравниваем хеш паролей если совпадают генерим и отдаем токены
func (a *AuthUseCase) Login(ctx context.Context, req dto.LoginReq) (*dto.TokenPair, error) {
	MetaInfo := &models.MetaInfo{}
	var err error
	switch req.IsOrg {
	case false: // user
		MetaInfo, err = a.user.UserGetMetaInfo(ctx, req.Email)
	case true: // org
		MetaInfo, err = a.org.OrgGetMetaInfo(ctx, req.Email)
	}
	if err != nil {
		a.Logger.Error(
			"failed login account",
			zap.Error(err),
		)
		return nil, err
	}
	// если не активирован
	if !MetaInfo.Verified {
		// то проверяем не стух ли он еще
		if validation.IsAccountExpired(MetaInfo.CreatedAt) {
			return nil, ErrAccountExpired
		}
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
// Отправляем на указанную почту код подтверждения
func (a *AuthUseCase) UserRegister(ctx context.Context, req dto.UserRegisterReq) (int, error) {
	_, err := a.user.UserIsExist(ctx, req.Email)
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("UserIsExist", err.Error()),
		)
		return 0, err
	}
	hash, err := passwd.GetHash(req.Password)
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("GetHash", err.Error()),
		)
		return 0, err
	}
	creds := entity.Credentials{
		Login:      req.Credentials.Email,
		PasswdHash: hash,
	}
	// Создали юзера
	userID, err := a.user.UserSave(ctx, &req.UserInfo, &creds)
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("UserSave", err.Error()),
		)
		return 0, err
	}
	// Генерируем код
	code, err := verification.GenerateCode()
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("GenerateCode", err.Error()),
		)
		return 0, err
	}
	// Сохраняем его в БД
	if err := a.user.UserSaveCode(ctx, code, userID); err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("UserSaveCode", err.Error()),
		)
		return 0, err
	}
	// Отправляем на почту
	if err := a.mail.SendVerifyCode(req.Email, code); err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("SendVerifyCode", err.Error()),
		)
		return 0, err
	}
	return userID, nil
}

// Аналогично работе с пользователем
func (a *AuthUseCase) OrgRegister(ctx context.Context, req dto.OrgRegisterReq) (int, error) {
	_, err := a.org.OrgIsExist(ctx, req.Email)
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("OrgIsExist", err.Error()),
		)
		return 0, err
	}
	hash, err := passwd.GetHash(req.Password)
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("GetHash", err.Error()),
		)
		return 0, err
	}
	creds := entity.Credentials{
		Login:      req.Credentials.Email,
		PasswdHash: hash,
	}
	orgID, err := a.org.OrgSave(ctx, &req.OrgInfo, &creds, req.City)
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("OrgSave", err.Error()),
		)
		return 0, err
	}
	// Генерируем код
	code, err := verification.GenerateCode()
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("GenerateCode", err.Error()),
		)
		return 0, err
	}
	// Сохраняем его в БД
	if err := a.org.OrgSaveCode(ctx, code, orgID); err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("OrgSaveCode", err.Error()),
		)
		return 0, err
	}
	// Отправляем на почту
	if err := a.mail.SendVerifyCode(req.Email, code); err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("SendVerifyCode", err.Error()),
		)
		return 0, err
	}
	return orgID, nil
}

func (a *AuthUseCase) SendCodeRetry(ctx context.Context, req dto.SendCodeReq) error {
	// Генерируем код
	code, err := verification.GenerateCode()
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("GenerateCode", err.Error()),
		)
		return err
	}
	switch req.IsOrg {
	case false:
		err = a.user.UserSaveCode(ctx, code, req.ID)
	case true:
		err = a.org.OrgSaveCode(ctx, code, req.ID)
	}
	// Сохраняем его в БД
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("User/OrgSaveCode", err.Error()),
		)
		return err
	}
	// Отправляем на почту
	if err := a.mail.SendVerifyCode(req.Email, code); err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("SendVerifyCode", err.Error()),
		)
		return err
	}
	return nil
}

func (a *AuthUseCase) VerifyCode(ctx context.Context, req dto.VerifyCodeReq) (*dto.TokenPair, error) {
	var exp time.Time
	var err error
	switch req.IsOrg {
	case false:
		exp, err = a.user.UserCode(ctx, req.Code, req.ID)
	case true:
		exp, err = a.org.OrgCode(ctx, req.Code, req.ID)
	}
	if err != nil {
		a.Logger.Error(
			"failed to verify code",
			zap.String("User/OrgCode", err.Error()),
		)
		return nil, err
	}
	if ok := validation.IsCodeExpired(exp); ok {
		a.Logger.Error(
			"failed to verify code",
			zap.String("IsCodeExpired", ErrCodeExpired.Error()),
		)
		return nil, ErrCodeExpired
	}
	// Активируем аккаунт
	switch req.IsOrg {
	case false:
		err = a.user.UserActivateAccount(ctx, req.ID)
	case true:
		err = a.org.OrgActivateAccount(ctx, req.ID)
	}
	if err != nil {
		a.Logger.Error(
			"failed to verify code",
			zap.String("UserActivateAccount", err.Error()),
		)
		return nil, err
	}
	// Генерим токен
	tokens, err := jwtlib.NewTokenPair(
		a.secret,
		a.TokenCfg,
		&entity.TokenMetadata{
			ID:    uint64(req.ID),
			IsOrg: req.IsOrg,
		},
	)
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("NewTokenPair", err.Error()),
		)
		return nil, err
	}
	return tokens, nil
}

func (a *AuthUseCase) UpdateAccessToken(ctx context.Context, req *jwt.Token) (*dto.AccessToken, error) {
	// Валидируем Claims токена. Есть ли они и нормальные ли.
	err := validation.ValidateTokenClaims(req)
	if err != nil {
		return nil, err
	}
	// Здесь уже спокойно кастую если выше проблем не возникло
	tmp := req.Claims.(jwt.MapClaims)["id"].(float64)
	id := uint64(tmp)
	metadata := &entity.TokenMetadata{
		ID:    id,
		IsOrg: req.Claims.(jwt.MapClaims)["is_org"].(bool),
	}
	token, err := jwtlib.NewToken(a.secret, a.TokenCfg, metadata, "access")
	if err != nil {
		return nil, err
	}
	return &dto.AccessToken{
		Token: token,
	}, nil
}
