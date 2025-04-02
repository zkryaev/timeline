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

type Middleware struct {
	secret *rsa.PrivateKey
	logger *zap.Logger
}

func New(key *rsa.PrivateKey, logger *zap.Logger) *Middleware {
	return &Middleware{
		secret: key,
		logger: logger,
	}
}

func (m *Middleware) HandlerLogs(next http.Handler) http.Handler {
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
		ctx := context.WithValue(r.Context(), "uuid", uuid.String()) //nolint:go-staticcheck // keep it simple
		start := time.Now()
		next.ServeHTTP(rw, r.WithContext(ctx))
		m.logger.Info("", zap.String("uuid", uuid.String()), zap.Int("code", rw.StatusCode()), zap.Duration("elapsed", time.Since(start)))
	})
}

func (m *Middleware) ExtractToken(r *http.Request) (*jwt.Token, error) {
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
			return nil, errors.Join(ErrInvalidSign, fmt.Errorf("token sign: %s", token.Method.Alg()))
		}
		publicKey := &m.secret.PublicKey
		return publicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))
	return token, err
}

// Валидацпия access токена
func (m *Middleware) IsAccessTokenValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.ExtractToken(r)

		if err != nil || !token.Valid {
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}
		if err = validation.ValidateTokenClaims(token.Claims); err != nil {
			http.Error(w, "token claims are invalid", http.StatusUnauthorized)
			return
		}
		if token.Claims.(jwt.MapClaims)["type"].(string) != "access" {
			http.Error(w, "access token wasn't provided", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
