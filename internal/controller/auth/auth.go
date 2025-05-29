package auth

import (
	"context"
	"errors"
	"net/http"
	"timeline/internal/controller/auth/middleware"
	"timeline/internal/controller/common"
	"timeline/internal/controller/scope"
	"timeline/internal/entity/dto/authdto"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthUseCase interface {
	Login(ctx context.Context, logger *zap.Logger, req *authdto.LoginReq) (*authdto.TokenPair, error)
	UserRegister(ctx context.Context, logger *zap.Logger, req *authdto.UserRegisterReq) (*authdto.TokenPair, error)
	OrgRegister(ctx context.Context, logger *zap.Logger, req *authdto.OrgRegisterReq) (*authdto.TokenPair, error)
	SendCodeRetry(ctx context.Context, logger *zap.Logger, req *authdto.SendCodeReq) error
	VerifyCode(ctx context.Context, logger *zap.Logger, req *authdto.VerifyCodeReq) error
	UpdateAccessToken(ctx context.Context, logger *zap.Logger, req *jwt.Token) (*authdto.AccessToken, error)
}

type AuthCtrl struct {
	usecase    AuthUseCase
	Logger     *zap.Logger
	Middleware middleware.Middleware
	settings   *scope.Settings
}

func New(usecase AuthUseCase, middleware middleware.Middleware, logger *zap.Logger, settings *scope.Settings) *AuthCtrl {
	return &AuthCtrl{
		usecase:    usecase,
		Logger:     logger,
		Middleware: middleware,
		settings:   settings,
	}
}

// @Summary Login
// @Description Authorizes a entity and returns a token pair
// @tags auth
// @Accept  json
// @Produce json
// @Param   request body authdto.LoginReq true " "
// @Success 200 {object} authdto.TokenPair
// @Failure 400
// @Failure 404
// @Failure 423
// @Failure 500
// @Router /auth/login [post]
func (a *AuthCtrl) Login(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(a.settings, a.Logger, r.Context())
	var req authdto.LoginReq
	if err := common.DecodeAndValidate(r, &req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := a.usecase.Login(r.Context(), logger, &req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrAccountNotFound):
			logger.Info("Login", zap.Error(err))
			http.Error(w, common.ErrAccountNotFound.Error(), http.StatusNotFound)
			return
		case errors.Is(err, common.ErrAccountExpired):
			logger.Info("Login", zap.Error(err))
			http.Error(w, common.ErrAccountExpired.Error(), http.StatusLocked)
			return
		default:
			logger.Error("Login", zap.Error(err))
			http.Error(w, common.ErrFailedLogin, http.StatusInternalServerError)
			return
		}
	}
	if err := common.WriteJSON(w, data); err != nil {
		logger.Error("WriteJSON", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// @Summary User registration
// @Description
// @tags auth
// @Accept  json
// @Produce  json
// @Param   request body authdto.UserRegisterReq true " "
// @Success 201 {object} authdto.TokenPair "User token"
// @Failure 400
// @Failure 500
// @Router /auth/registration/users [post]
func (a *AuthCtrl) UserRegister(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(a.settings, a.Logger, r.Context())
	var req authdto.UserRegisterReq
	if err := common.DecodeAndValidate(r, &req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := a.usecase.UserRegister(r.Context(), logger, &req)
	if err != nil {
		logger.Error("UserRegister", zap.Error(err))
		http.Error(w, common.ErrFailedRegister, http.StatusInternalServerError)
		return
	}
	if err := common.WriteJSON(w, data); err != nil {
		logger.Error("WriteJSON", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// @Summary Organization registration
// @Description
// @tags auth
// @Accept  json
// @Produce  json
// @Param   request body authdto.OrgRegisterReq true " "
// @Success 201 {object} authdto.TokenPair "Organization token"
// @Failure 400
// @Failure 500
// @Router /auth/registration/orgs [post]
func (a *AuthCtrl) OrganizationRegister(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(a.settings, a.Logger, r.Context())
	var req authdto.OrgRegisterReq
	if err := common.DecodeAndValidate(r, &req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := a.usecase.OrgRegister(r.Context(), logger, &req)
	if err != nil {
		logger.Error("OrgRegister", zap.Error(err))
		http.Error(w, common.ErrFailedRegister, http.StatusInternalServerError)
		return
	}
	if err := common.WriteJSON(w, data); err != nil {
		logger.Error("WriteJSON", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// @Summary Send code to email
// @Description
// @tags auth
// @Accept  json
// @Produce  json
// @Param   request body authdto.SendCodeReq true " "
// @Success 201
// @Failure 400
// @Failure 500
// @Router /auth/codes [post]
func (a *AuthCtrl) CodeSend(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(a.settings, a.Logger, r.Context())
	var req authdto.SendCodeReq
	if err := common.DecodeAndValidate(r, &req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := a.usecase.SendCodeRetry(r.Context(), logger, &req); err != nil {
		logger.Error("SendCodeRetry", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// @Summary Confirm code
// @Description Confirm code that has been sent to email
// @tags auth
// @Accept  json
// @Produce  json
// @Param   request body authdto.VerifyCodeReq true " "
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 410
// @Failure 500
// @Router /auth/codes [put]
func (a *AuthCtrl) CodeConfirm(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(a.settings, a.Logger, r.Context())
	var req authdto.VerifyCodeReq
	if err := common.DecodeAndValidate(r, &req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := a.usecase.VerifyCode(r.Context(), logger, &req); err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("VerifyCode", zap.Error(err))
			http.Error(w, common.ErrNotFound.Error(), http.StatusNotFound)
			return
		case errors.Is(err, common.ErrCodeExpired):
			logger.Info("VerifyCode", zap.Error(err))
			http.Error(w, common.ErrCodeExpired.Error(), http.StatusGone)
			return
		default:
			logger.Error("VerifyCode", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Refresh access token
// @Description
// @tags auth
// @Accept  json
// @Produce  json
// @Param   refresh_token header string true "Refresh Token"
// @Success 200 {object} authdto.AccessToken "Updated Access Token"
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /auth/token [put]
func (a *AuthCtrl) PutAccessToken(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(a.settings, a.Logger, r.Context())
	token, err := a.Middleware.ExtractToken(r)
	if err != nil {
		logger.Error("ExtractToken", zap.String("ExtractToken", err.Error()))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if tokenType := token.Claims.(jwt.MapClaims)["type"].(string); tokenType != "refresh" {
		logger.Info("invalid token", zap.String("token_type", tokenType))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	refreshedAccessToken, err := a.usecase.UpdateAccessToken(r.Context(), logger, token)
	if err != nil {
		logger.Error("UpdateAccessToken", zap.Error(err))
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	if common.WriteJSON(w, refreshedAccessToken) != nil {
		logger.Error("WriteJSON", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
