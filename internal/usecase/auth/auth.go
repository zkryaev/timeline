package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"strings"
	"timeline/internal/config"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/infrastructure"
	"timeline/internal/infrastructure/mail"
	"timeline/internal/infrastructure/mapper/codemap"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/infrastructure/mapper/usermap"
	"timeline/internal/infrastructure/models"
	jwtlib "timeline/internal/libs/jwtlib"
	"timeline/internal/libs/passwd"
	"timeline/internal/libs/verification"
	"timeline/internal/usecase/auth/validation"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountExpired     = errors.New("account expired")
	ErrCodeExpired        = errors.New("code expired")
)

type AuthUseCase struct {
	secret   *rsa.PrivateKey
	user     infrastructure.UserRepository
	org      infrastructure.OrgRepository
	code     infrastructure.CodeRepository
	mail     infrastructure.Mail
	TokenCfg config.Token
	Logger   *zap.Logger
}

func New(key *rsa.PrivateKey, userRepo infrastructure.UserRepository, orgRepo infrastructure.OrgRepository, codeRepo infrastructure.CodeRepository, mailSrv infrastructure.Mail, cfg config.Token, logger *zap.Logger) *AuthUseCase {
	return &AuthUseCase{
		secret:   key,
		user:     userRepo,
		org:      orgRepo,
		code:     codeRepo,
		mail:     mailSrv,
		TokenCfg: cfg,
		Logger:   logger,
	}
}

func (a *AuthUseCase) Login(ctx context.Context, req *authdto.LoginReq) (*authdto.TokenPair, error) {
	exp, err := a.code.AccountExpiration(ctx, req.Email, req.IsOrg)
	if err != nil {
		a.Logger.Error(
			"failed login account",
			zap.String("AccountExpiration", err.Error()),
		)
		return nil, err
	}
	// если не активирован
	if !exp.Verified {
		// то проверяем не стух ли он еще
		if validation.IsAccountExpired(exp.CreatedAt) {
			return nil, ErrAccountExpired
		}
	}
	if err := passwd.CompareWithHash(req.Password, exp.Hash); err != nil {
		a.Logger.Error(
			"failed login account",
			zap.String("CompareWithHash", err.Error()),
		)
		return nil, err
	}

	tokens, err := jwtlib.NewTokenPair(
		a.secret,
		a.TokenCfg,
		&entity.TokenMetadata{
			ID:    uint64(exp.ID),
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

func (a *AuthUseCase) UserRegister(ctx context.Context, req *authdto.UserRegisterReq) (*authdto.RegisterResp, error) {
	hash, err := passwd.GetHash(req.Password)
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("GetHash", err.Error()),
		)
		return nil, err
	}
	req.Credentials.Password = hash
	id, err := uuid.NewRandom()
	req.UUID = id.String()
	userID, err := a.user.UserSave(ctx, usermap.UserRegisterToModel(req))
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("UserSave", err.Error()),
		)
		return nil, err
	}
	// Генерируем код
	code, err := verification.GenerateCode()
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("GenerateCode", err.Error()),
		)
		return nil, err
	}
	a.mail.SendMsg(&models.Message{
		Email: req.Email,
		Type:  mail.VerificationType,
		Value: code,
	})
	return &authdto.RegisterResp{Id: userID}, nil
}

func (a *AuthUseCase) OrgRegister(ctx context.Context, req *authdto.OrgRegisterReq) (*authdto.RegisterResp, error) {
	hash, err := passwd.GetHash(req.Password)
	if err != nil {
		a.Logger.Error(
			"failed to register org",
			zap.String("GetHash", err.Error()),
		)
		return nil, err
	}
	req.Credentials.Password = hash
	req.Name = strings.ToLower(req.Name)
	id, err := uuid.NewRandom()
	req.UUID = id.String()
	orgID, err := a.org.OrgSave(ctx, orgmap.RegisterReqToModel(req))
	if err != nil {
		a.Logger.Error(
			"failed to register org",
			zap.String("OrgSave", err.Error()),
		)
		return nil, err
	}

	code, err := verification.GenerateCode()
	if err != nil {
		a.Logger.Error(
			"failed to register org",
			zap.String("GenerateCode", err.Error()),
		)
		return nil, err
	}
	a.mail.SendMsg(&models.Message{
		Email: req.Email,
		Type:  mail.VerificationType,
		Value: code,
	})
	return &authdto.RegisterResp{Id: orgID}, nil
}

func (a *AuthUseCase) SendCodeRetry(ctx context.Context, req *authdto.SendCodeReq) {
	code, err := verification.GenerateCode()
	if err != nil {
		a.Logger.Error(
			"retry send code failed",
			zap.String("GenerateCode", err.Error()),
		)
		return
	}
	a.mail.SendMsg(&models.Message{
		Email: req.Email,
		Type:  mail.VerificationType,
		Value: code,
	})
}

func (a *AuthUseCase) VerifyCode(ctx context.Context, req *authdto.VerifyCodeReq) (*authdto.TokenPair, error) {
	exp, err := a.code.VerifyCode(ctx, codemap.ToModel(req))
	if err != nil {
		a.Logger.Error(
			"failed to verify code",
			zap.String("VerifyCode", err.Error()),
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

	if err = a.code.ActivateAccount(ctx, req.ID, req.IsOrg); err != nil {
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

func (a *AuthUseCase) UpdateAccessToken(ctx context.Context, req *jwt.Token) (*authdto.AccessToken, error) {
	// Валидируем Claims токена. Есть ли они и нормальные ли.
	err := validation.ValidateTokenClaims(req.Claims)
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
	return &authdto.AccessToken{
		Token: token,
	}, nil
}
