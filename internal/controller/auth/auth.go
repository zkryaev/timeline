package auth

import (
	"context"
	"net/http"
	"timeline/internal/entity/dto"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type Auth interface {
	Login(ctx context.Context, req dto.LoginReq) (*dto.TokenPair, error)
	UserRegister(ctx context.Context, req dto.UserRegisterReq) (int, error)
	OrgRegister(ctx context.Context, req dto.OrgRegisterReq) (int, error)
	SendCodeRetry(ctx context.Context, req dto.SendCodeReq) error
	VerifyCode(ctx context.Context, req dto.VerifyCodeReq) (*dto.TokenPair, error)
	UpdateAccessToken(ctx context.Context, req *jwt.Token) (string, error)
}

type Middleware interface {
	ExtractToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, error)
	IsTokenValid(next http.Handler) http.Handler
}

type AuthCtrl struct {
	usecase    Auth
	Logger     *zap.Logger
	Middleware Middleware
	json       jsoniter.API
	validator  validator.Validate
}

func New(usecase Auth, middleware Middleware, logger *zap.Logger, jsoniter jsoniter.API, validator validator.Validate) *AuthCtrl {
	return &AuthCtrl{
		usecase:    usecase,
		Logger:     logger,
		Middleware: middleware,
		json:       jsoniter,
		validator:  validator,
	}
}

// Авторизация объекта. Принимает логин + пароль. Возвращает пару токенов
func (a *AuthCtrl) Login(w http.ResponseWriter, r *http.Request) {
	// декодируем json
	var req dto.LoginReq
	if a.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	// валидация полей
	if a.validator.Struct(&req) != nil {
		http.Error(w, "Data is not valid", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	data, err := a.usecase.Login(ctx, req)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}
	// отдаем токен
	if a.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Регистрация пользователя. Получает информацию о пользователе. Возвращает пару токенов
func (a *AuthCtrl) UserRegister(w http.ResponseWriter, r *http.Request) {
	// декодируем json
	var req dto.UserRegisterReq
	if a.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusBadRequest)
		return
	}
	// валидация полей
	if a.validator.Struct(&req) != nil {
		http.Error(w, "Data is not valid", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	data, err := a.usecase.UserRegister(ctx, req)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	// отдаем токен
	if a.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Регистрация организации. Получает информацию о организации. Возвращает пару токенов
func (a *AuthCtrl) OrgRegister(w http.ResponseWriter, r *http.Request) {
	// декодируем json
	var req dto.OrgRegisterReq
	if a.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusBadRequest)
		return
	}
	// валидация полей
	if a.validator.Struct(&req) != nil {
		http.Error(w, "Data is not valid", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	data, err := a.usecase.OrgRegister(ctx, req)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	// отдаем токен
	if a.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *AuthCtrl) SendCodeRetry(w http.ResponseWriter, r *http.Request) {
	// декодируем json
	var req dto.SendCodeReq
	if a.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusBadRequest)
		return
	}
	// валидация полей
	if a.validator.Struct(&req) != nil {
		http.Error(w, "Data is not valid", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := a.usecase.SendCodeRetry(ctx, req); err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *AuthCtrl) VerifyCode(w http.ResponseWriter, r *http.Request) {
	// декодируем json
	var req dto.VerifyCodeReq
	if a.json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusBadRequest)
		return
	}
	// валидация полей
	if a.validator.Struct(&req) != nil {
		http.Error(w, "Data is not valid", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	data, err := a.usecase.VerifyCode(ctx, req)
	if err != nil {
		http.Error(w, "Code or account expired", http.StatusBadRequest)
		return
	}
	// отдаем токен
	if a.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Принимает refresh токен. Возвращает новый access токен
func (a *AuthCtrl) UpdateAccessToken(w http.ResponseWriter, r *http.Request) {
	token, err := a.Middleware.ExtractToken(w, r)
	if err != nil {
		a.Logger.Error(
			"failed update access token",
			zap.String("ExtractToken", err.Error()),
		)
		http.Error(w, "Invalid token field", http.StatusBadRequest)
		return
	}
	if token.Claims.(jwt.MapClaims)["type"] != "refresh" {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}
	// TODO: Нет работы с контекстами
	ctx := context.Background()
	refreshedAccessToken, err := a.usecase.UpdateAccessToken(ctx, token)
	if err != nil {
		a.Logger.Error(
			"failed update access token",
			zap.String("usecase.UpdateAccessToken", err.Error()),
		)
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]string{
		"access_token": refreshedAccessToken,
	}

	// Отправляем JSON-ответ
	if err := a.json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
