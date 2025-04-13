package middleware

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"timeline/internal/controller/common"
	"timeline/internal/controller/settings"
	"timeline/internal/usecase/auth/validation"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrTokenNotFound   = errors.New("token not found")
	ErrAuthHeaderEmpty = errors.New("auth header empty")
	ErrInvalidSign     = errors.New("invalid sign method")
)

type UUID string
type TokenData string

type Middleware interface {
	ExtractToken(r *http.Request) (*jwt.Token, error)
	Authorization(next http.Handler) http.Handler
	HandlerLogs(next http.Handler) http.Handler
}

type middleware struct {
	secret *rsa.PrivateKey
	logger *zap.Logger
	routes settings.Routes
}

func New(key *rsa.PrivateKey, logger *zap.Logger, routes settings.Routes) Middleware {
	return &middleware{
		secret: key,
		logger: logger,
		routes: routes,
	}
}

func (m *middleware) HandlerLogs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &common.ResponseWriter{
			ResponseWriter: w,
		}
		uuid, err := uuid.NewRandom()
		if err != nil {
			m.logger.Error("HandlerLogs", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		uri, err := url.QueryUnescape(r.RequestURI)
		if err != nil {
			uri = r.RequestURI // Если декодирование не удалось, используем оригинальный URI
		}
		m.logger.Info(r.Method, zap.String("uuid", uuid.String()), zap.String("uri", uri))
		ctx := context.WithValue(r.Context(), UUID("uuid"), uuid.String()) // nolint:go-staticcheck // keep it simple
		start := time.Now()
		next.ServeHTTP(rw, r.WithContext(ctx))
		m.logger.Info("", zap.String("uuid", uuid.String()), zap.Int("code", rw.StatusCode()), zap.Duration("elapsed", time.Since(start)))
	})
}

func (m *middleware) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.ExtractToken(r)
		if err != nil {
			m.logger.Error("ExtractToken", zap.Error(err))
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		if err := validation.IsTokenValid(token); err != nil {
			m.logger.Error("IsTokenValid", zap.Error(err))
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		tdata := common.GetTokenData(token.Claims)
		if err := m.routes.HasAccess(tdata, r.RequestURI, r.Method); err != nil {
			m.logger.Error("HasAccess", zap.Error(err))
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), TokenData("token"), tdata) // nolint:go-staticcheck // keep it simple
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *middleware) ExtractToken(r *http.Request) (*jwt.Token, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, ErrAuthHeaderEmpty
	}
	tokenString := strings.Split(authHeader, " ")[1]
	if tokenString == "" {
		return nil, ErrTokenNotFound
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("%s: %s", ErrInvalidSign, token.Method.Alg())
		}
		publicKey := &m.secret.PublicKey
		return publicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))
	return token, err
}
