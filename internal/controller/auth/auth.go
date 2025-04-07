package auth

import (
	"context"
	"errors"
	"net/http"
	"timeline/internal/controller/auth/middleware"
	"timeline/internal/controller/common"
	"timeline/internal/entity/dto/authdto"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthUseCase interface {
	Login(ctx context.Context, logger *zap.Logger, req *authdto.LoginReq) (*authdto.TokenPair, error)
	UserRegister(ctx context.Context, logger *zap.Logger, req *authdto.UserRegisterReq) (*authdto.RegisterResp, error)
	OrgRegister(ctx context.Context, logger *zap.Logger, req *authdto.OrgRegisterReq) (*authdto.RegisterResp, error)
	SendCodeRetry(ctx context.Context, logger *zap.Logger, req *authdto.SendCodeReq)
	VerifyCode(ctx context.Context, logger *zap.Logger, req *authdto.VerifyCodeReq) (*authdto.TokenPair, error)
	UpdateAccessToken(ctx context.Context, logger *zap.Logger, req *jwt.Token) (*authdto.AccessToken, error)
}

type AuthCtrl struct {
	usecase    AuthUseCase
	Logger     *zap.Logger
	Middleware middleware.Middleware
}

func New(usecase AuthUseCase, middleware middleware.Middleware, logger *zap.Logger) *AuthCtrl {
	return &AuthCtrl{
		usecase:    usecase,
		Logger:     logger,
		Middleware: middleware,
	}
}

// @Summary Login
// @Description Authorizes a user and returns a token pair
// @Tags Auth
// @Accept  json
// @Produce json
// @Param   request body authdto.LoginReq true "Login Request"
// @Success 200 {object} authdto.TokenPair
// @Failure 400
// @Failure 404
// @Failure 423 {string} string "account expired
// @Failure 500
// @Router /auth/login [post]
func (a *AuthCtrl) Login(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := a.Logger.With(zap.String("uuid", uuid))
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

// @Summary Register User
// @Description Registers a user and returns a token pair
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   request body authdto.UserRegisterReq true "User Register Request"
// @Success 201 {object} authdto.RegisterResp "User ID"
// @Failure 400
// @Failure 500
// @Router /auth/users [post]
func (a *AuthCtrl) UserRegister(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := a.Logger.With(zap.String("uuid", uuid))
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

// @Summary Register Organization
// @Description Registers an organization and returns a token pair
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   request body authdto.OrgRegisterReq true "Organization Register Request"
// @Success 201 {object} authdto.RegisterResp "Organization ID"
// @Failure 400
// @Failure 500
// @Router /auth/orgs [post]
func (a *AuthCtrl) OrgRegister(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := a.Logger.With(zap.String("uuid", uuid))
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

// @Summary Send Code Retry
// @Description Sends a code retry request
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   request body authdto.SendCodeReq true "Send Code Request"
// @Success 201 {string} string "Code resent successfully"
// @Failure 400
// @Failure 500
// @Router /auth/codes/send [post]
func (a *AuthCtrl) SendCodeRetry(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := a.Logger.With(zap.String("uuid", uuid))
	var req authdto.SendCodeReq
	if err := common.DecodeAndValidate(r, &req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	a.usecase.SendCodeRetry(r.Context(), logger, &req)
	w.WriteHeader(http.StatusCreated)
}

// @Summary Verify Code
// @Description Verifies the code and returns a token pair
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   request body authdto.VerifyCodeReq true "Verify Code Request"
// @Success 200 {object} authdto.TokenPair
// @Failure 400
// @Failure 404
// @Failure 410 {string} string "code expired"
// @Failure 500
// @Router /auth/codes/verify [post]
func (a *AuthCtrl) VerifyCode(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := a.Logger.With(zap.String("uuid", uuid))
	var req authdto.VerifyCodeReq
	if err := common.DecodeAndValidate(r, &req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := a.usecase.VerifyCode(r.Context(), logger, &req)
	if err != nil {
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
			http.Error(w, common.ErrFailedRegister, http.StatusInternalServerError)
			return
		}
	}
	if err := common.WriteJSON(w, data); err != nil {
		logger.Error("WriteJSON", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// @Summary Update Access Token
// @Description Updates the access token using a refresh token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   refresh_token header string true "Refresh Token"
// @Success 200 {object} authdto.AccessToken "New Access Token"
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /auth/tokens/refresh [put]
func (a *AuthCtrl) UpdateAccessToken(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := a.Logger.With(zap.String("uuid", uuid))
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
