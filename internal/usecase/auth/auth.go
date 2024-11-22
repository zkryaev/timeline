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
	"timeline/internal/repository"
	"timeline/internal/repository/mail/notify"
	"timeline/internal/repository/mapper/codemap"
	"timeline/internal/repository/mapper/orgmap"
	"timeline/internal/repository/mapper/usermap"
	"timeline/internal/usecase/auth/validation"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountExpired     = errors.New("account expired")
	ErrCodeExpired        = errors.New("code expired")
)

var (
	// Получено экспериментальным путем.
	// Верхняя замеренная граница отправки 10 секунд, в случае сервиса gmail.com
	SendEmailTimeout = 10 * time.Second
)

type AuthUseCase struct {
	secret   *rsa.PrivateKey
	user     repository.UserRepository
	org      repository.OrgRepository
	code     repository.CodeRepository
	mail     notify.Mail
	TokenCfg config.Token
	Logger   *zap.Logger
}

func New(key *rsa.PrivateKey, userRepo repository.UserRepository, orgRepo repository.OrgRepository, codeRepo repository.CodeRepository, mailSrv notify.Mail, cfg config.Token, logger *zap.Logger) *AuthUseCase {
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

func (a *AuthUseCase) Login(ctx context.Context, req *dto.LoginReq) (*dto.TokenPair, error) {
	// TODO: Пока так, а вообще CRON по идее будет чистить сгоревшие аккаунты
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

func (a *AuthUseCase) UserRegister(ctx context.Context, req *dto.UserRegisterReq) (*dto.RegisterResp, error) {
	hash, err := passwd.GetHash(req.Password)
	if err != nil {
		a.Logger.Error(
			"failed to register user",
			zap.String("GetHash", err.Error()),
		)
		return nil, err
	}
	req.Credentials.Password = hash
	// Создали юзера
	userID, err := a.user.UserSave(ctx, usermap.ToModel(req))
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
	// TODO: подумать над этим еще
	ctxOnSend, _ := context.WithTimeout(context.Background(), SendEmailTimeout)
	go a.codeProccessing(ctxOnSend, &dto.VerifyCodeReq{
		ID:    userID,
		Email: req.Email,
		Code:  code,
		IsOrg: false,
	})
	return &dto.RegisterResp{Id: userID}, nil
}

// (Experimental)
// Выполняет отправку кода на почту и сохранение кода в БД.
// Чтобы не флудить логгами, отображаются только Error и Warn.
// Error - явная ошибка либо в используемом сервисе почты, либо в БД.
// Warn - истечение таймаута = неизвестная, неявная ошибка.
func (a *AuthUseCase) codeProccessing(ctx context.Context, metaInfo *dto.VerifyCodeReq) {
	maxRetries := 2
	retryDelay := 500 * time.Millisecond
	done := make(chan error, 1)
	defer close(done)
	go func() {
		var err error
		for try := maxRetries; try > 0; try-- {
			// start := time.Now()
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			default:
				err = a.mail.SendVerifyCode(metaInfo.Email, metaInfo.Code)
			}
			// log.Println("Send code to email took", time.Now().Sub(start), try)
			if err == nil {
				break
			}
			time.Sleep(retryDelay)
		} // TODO: Что насчет error groups?
		// start := time.Now()
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		default:
			if err != nil {
				done <- err
				return
			} else if err = a.code.SaveVerifyCode(ctx, codemap.ToModel(metaInfo)); err == nil {
				// log.Println("Save code to DB took", time.Now().Sub(start))
				done <- nil
				return
			} else {
				done <- err
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		a.Logger.Warn(
			"processing code failed due to timeout",
			zap.String("email", metaInfo.Email),
		)
		return
	case err := <-done:
		if err != nil {
			a.Logger.Error("failed to send code to email", zap.Error(err))
		}
		return
	}
}

// Аналогично работе с пользователем
func (a *AuthUseCase) OrgRegister(ctx context.Context, req *dto.OrgRegisterReq) (*dto.RegisterResp, error) {
	hash, err := passwd.GetHash(req.Password)
	if err != nil {
		a.Logger.Error(
			"failed to register org",
			zap.String("GetHash", err.Error()),
		)
		return nil, err
	}
	req.Credentials.Password = hash
	orgID, err := a.org.OrgSave(ctx, orgmap.ToModel(req))
	if err != nil {
		a.Logger.Error(
			"failed to register org",
			zap.String("OrgSave", err.Error()),
		)
		return nil, err
	}
	// Генерируем код
	code, err := verification.GenerateCode()
	if err != nil {
		a.Logger.Error(
			"failed to register org",
			zap.String("GenerateCode", err.Error()),
		)
		return nil, err
	}
	ctxOnSend, _ := context.WithTimeout(context.Background(), SendEmailTimeout)
	go a.codeProccessing(ctxOnSend, &dto.VerifyCodeReq{
		ID:    orgID,
		Email: req.Email,
		Code:  code,
		IsOrg: true,
	})
	return &dto.RegisterResp{Id: orgID}, nil
}

func (a *AuthUseCase) SendCodeRetry(ctx context.Context, req *dto.SendCodeReq) {
	// Генерируем код
	code, err := verification.GenerateCode()
	if err != nil {
		a.Logger.Error(
			"retry send code failed",
			zap.String("GenerateCode", err.Error()),
		)
		return
	}
	ctxOnSend, _ := context.WithTimeout(context.Background(), 4*time.Second)
	a.codeProccessing(ctxOnSend, &dto.VerifyCodeReq{
		ID:    req.ID,
		Email: req.Email,
		IsOrg: req.IsOrg,
		Code:  code,
	})
}

func (a *AuthUseCase) VerifyCode(ctx context.Context, req *dto.VerifyCodeReq) (*dto.TokenPair, error) {
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
