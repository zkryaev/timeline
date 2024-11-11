package auth

import (
	"context"
	"net/http"
	"timeline/internal/model/dto"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type Auth interface {
	Login(ctx context.Context, req dto.LoginReq) (*dto.TokenPair, error)
	UserRegister(ctx context.Context, req dto.UserRegisterReq) (*dto.TokenPair, error)
	OrgRegister(ctx context.Context, req dto.OrgRegisterReq) (*dto.TokenPair, error)
	UpdateRefreshToken(ctx context.Context, req *jwt.Token) (string, error)
}

type Middleware interface {
	ExtractToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, error)
	IsTokenValid(next http.Handler) http.Handler
	IsRefreshToken(next http.Handler) http.Handler
}

type AuthCtrl struct {
	Usecase    Auth
	Logger     *zap.Logger
	Middleware Middleware
	json       jsoniter.API
	validator  validator.Validate
}

func New(usecase Auth, middleware Middleware, logger *zap.Logger, jsoniter jsoniter.API, validator validator.Validate) *AuthCtrl {
	return &AuthCtrl{
		Usecase:    usecase,
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
	data, err := a.Usecase.Login(ctx, req)
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
	data, err := a.Usecase.UserRegister(ctx, req)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}
	// отдаем токен
	if a.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	data, err := a.Usecase.OrgRegister(ctx, req)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
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
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	ctx := context.Background()
	refreshedAccessToken, err := a.Usecase.UpdateRefreshToken(ctx, token)
	if err != nil {
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
