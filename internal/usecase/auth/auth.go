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
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/infrastructure/mail"
	"timeline/internal/infrastructure/mapper/codemap"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/infrastructure/mapper/usermap"
	"timeline/internal/infrastructure/models"
	jwtlib "timeline/internal/sugar/jwtlib"
	"timeline/internal/sugar/passwd"
	"timeline/internal/sugar/verification"
	"timeline/internal/usecase/auth/validation"
	"timeline/internal/usecase/common"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountExpired     = errors.New("account expired")
	ErrCodeExpired        = errors.New("code expired")
	ErrAccountNotFound    = postgres.ErrAccountNotFound
)

type AuthUseCase struct {
	secret   *rsa.PrivateKey
	user     infrastructure.UserRepository
	org      infrastructure.OrgRepository
	code     infrastructure.CodeRepository
	mail     infrastructure.Mail
	TokenCfg config.Token
}

func New(key *rsa.PrivateKey, userRepo infrastructure.UserRepository, orgRepo infrastructure.OrgRepository, codeRepo infrastructure.CodeRepository, mailSrv infrastructure.Mail, cfg config.Token) *AuthUseCase {
	return &AuthUseCase{
		secret:   key,
		user:     userRepo,
		org:      orgRepo,
		code:     codeRepo,
		mail:     mailSrv,
		TokenCfg: cfg,
	}
}

func (a *AuthUseCase) Login(ctx context.Context, logger *zap.Logger, req *authdto.LoginReq) (*authdto.TokenPair, error) {
	exp, err := a.code.AccountExpiration(ctx, req.Email, req.IsOrg)
	if err != nil {
		if errors.Is(err, postgres.ErrAccountNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	logger.Info("Fetched account metadata from DB")
	if !exp.Verified { // is activated
		if validation.IsAccountExpired(exp.CreatedAt) {
			return nil, ErrAccountExpired
		}
	}
	if err = passwd.CompareWithHash(req.Password, exp.Hash); err != nil {
		return nil, err
	}
	logger.Info("Password is correct")
	tokens, err := jwtlib.NewTokenPair(
		a.secret,
		a.TokenCfg,
		&entity.TokenData{
			ID:    exp.ID,
			IsOrg: req.IsOrg,
		},
	)
	if err != nil {
		return nil, err
	}
	logger.Info("Token pair have been generated")
	return tokens, nil
}

func (a *AuthUseCase) UserRegister(ctx context.Context, logger *zap.Logger, req *authdto.UserRegisterReq) (*authdto.RegisterResp, error) {
	hash, err := passwd.GetHash(req.Password)
	if err != nil {
		return nil, err
	}
	req.Credentials.Password = hash
	id, _ := uuid.NewRandom()
	req.UUID = id.String()
	userID, err := a.user.UserSave(ctx, usermap.UserRegisterToModel(req))
	if err != nil {
		return nil, err
	}
	logger.Info("User has been saved to DB")
	code, err := verification.GenerateCode(logger)
	if err != nil {
		return nil, err
	}
	info := &models.CodeInfo{
		ID:    userID,
		Code:  code,
		IsOrg: false,
	}
	if err := a.code.SaveVerifyCode(ctx, info); err != nil {
		return nil, err
	}
	logger.Info("Code has been saved to DB")
	a.mail.SendMsg(&models.Message{
		Email: req.Email,
		Type:  mail.VerificationType,
		Value: code,
	})
	logger.Info("Verification code has been sent to user's email")
	return &authdto.RegisterResp{ID: userID}, nil
}

func (a *AuthUseCase) OrgRegister(ctx context.Context, logger *zap.Logger, req *authdto.OrgRegisterReq) (*authdto.RegisterResp, error) {
	hash, err := passwd.GetHash(req.Password)
	if err != nil {
		return nil, err
	}
	req.Credentials.Password = hash
	req.Name = strings.ToLower(req.Name)
	id, _ := uuid.NewRandom()
	req.UUID = id.String()
	orgID, err := a.org.OrgSave(ctx, orgmap.RegisterReqToModel(req))
	if err != nil {
		return nil, err
	}
	logger.Info("Org has been saved to DB")
	code, err := verification.GenerateCode(logger)
	if err != nil {
		return nil, err
	}
	logger.Info("Verification code has been generated")
	info := &models.CodeInfo{
		ID:    orgID,
		Code:  code,
		IsOrg: true,
	}
	if err := a.code.SaveVerifyCode(ctx, info); err != nil {
		return nil, err
	}
	logger.Info("Code has been saved to DB")
	a.mail.SendMsg(&models.Message{
		Email: req.Email,
		Type:  mail.VerificationType,
		Value: code,
	})
	logger.Info("Verification code has been sent to user's email")
	return &authdto.RegisterResp{ID: orgID}, nil
}

func (a *AuthUseCase) SendCodeRetry(ctx context.Context, logger *zap.Logger, req *authdto.SendCodeReq) error {
	code, err := verification.GenerateCode(logger)
	if err != nil {
		return err
	}
	logger.Info("Verification code has been generated")
	info := &models.CodeInfo{
		ID:    req.ID,
		Code:  code,
		IsOrg: req.IsOrg,
	}
	if err := a.code.SaveVerifyCode(ctx, info); err != nil {
		return err
	}
	logger.Info("Code has been saved to DB")
	a.mail.SendMsg(&models.Message{
		Email: req.Email,
		Type:  mail.VerificationType,
		Value: code,
	})
	logger.Info("Verification code has been sent to user's email")
	return nil
}

func (a *AuthUseCase) VerifyCode(ctx context.Context, logger *zap.Logger, req *authdto.VerifyCodeReq) (*authdto.TokenPair, error) {
	exp, err := a.code.VerifyCode(ctx, codemap.ToModel(req))
	if err != nil {
		if errors.Is(err, postgres.ErrCodeNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Code was found")
	if ok := validation.IsCodeExpired(exp); ok {
		return nil, ErrCodeExpired
	}
	logger.Info("Code is fresh")
	if err = a.code.ActivateAccount(ctx, req.ID, req.IsOrg); err != nil {
		return nil, err
	}
	logger.Info("Account is activated")
	tokens, err := jwtlib.NewTokenPair(
		a.secret,
		a.TokenCfg,
		&entity.TokenData{
			ID:    req.ID,
			IsOrg: req.IsOrg,
		},
	)
	if err != nil {
		return nil, err
	}
	logger.Info("Token pair have been generated")
	return tokens, nil
}

func (a *AuthUseCase) UpdateAccessToken(_ context.Context, logger *zap.Logger, req *jwt.Token) (*authdto.AccessToken, error) {
	err := validation.ValidateTokenClaims(req.Claims)
	if err != nil {
		return nil, err
	}
	logger.Info("Token claims are correct")
	tmp := req.Claims.(jwt.MapClaims)["id"].(float64)
	id := int(tmp)
	metadata := &entity.TokenData{
		ID:    id,
		IsOrg: req.Claims.(jwt.MapClaims)["is_org"].(bool),
	}
	token, err := jwtlib.NewToken(a.secret, a.TokenCfg, metadata, "access")
	if err != nil {
		return nil, err
	}
	logger.Info("Access token has been generated")
	return &authdto.AccessToken{
		Token: token,
	}, nil
}
