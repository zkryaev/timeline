package middleware

import (
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
		start := time.Now()
		rw := common.NewResponseWriter(w)

		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		decodedURI, err := url.QueryUnescape(r.RequestURI)
		if err != nil {
			decodedURI = r.RequestURI // Если декодирование не удалось, используем оригинальный URI
		}
		logsText := fmt.Sprintf("method: %q, uri: %q, code: \"%d\", elapsed: %q", r.Method, decodedURI, rw.StatusCode(), duration)
		switch {
		case rw.StatusCode() >= http.StatusBadRequest && rw.StatusCode() < http.StatusInternalServerError:
			m.logger.Warn(logsText)
		case rw.StatusCode() >= http.StatusInternalServerError:
			m.logger.Error(logsText)
		default:
			m.logger.Info(logsText)
		}
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
