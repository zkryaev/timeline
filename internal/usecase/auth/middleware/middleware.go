package middleware

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type respWriterCustom struct {
	http.ResponseWriter
	statusCode int
	header     http.Header
}

func (rw *respWriterCustom) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *respWriterCustom) Write(b []byte) (int, error) {
	if rw.statusCode == 0 { // Если WriteHeader не был вызван
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *respWriterCustom) Header() http.Header {
	return rw.ResponseWriter.Header()
}

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

var (
	ErrTokenNotFound   = errors.New("token not found")
	ErrAuthHeaderEmpty = errors.New("auth header empty")
)

func (m *Middleware) HandlerLogs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &respWriterCustom{ResponseWriter: w}

		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		decodedURI, err := url.QueryUnescape(r.RequestURI)
		if err != nil {
			decodedURI = r.RequestURI // Если декодирование не удалось, используем оригинальный URI
		}
		logsText := fmt.Sprintf("method: %q, uri: %q, code: \"%d\", elapsed: %q", r.Method, decodedURI, rw.statusCode, duration)
		switch {
		case rw.statusCode >= http.StatusBadRequest && rw.statusCode < http.StatusInternalServerError:
			m.logger.Warn(logsText)
		case rw.statusCode >= http.StatusInternalServerError:
			m.logger.Error(logsText)
		default:
			m.logger.Info(logsText)
		}
	})
}

func (m *Middleware) ExtractToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, ErrAuthHeaderEmpty
	}
	tokenString := strings.Split(authHeader, " ")[1]
	if tokenString == "" {
		return nil, ErrTokenNotFound
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("неверный метод подписи: %v", token.Header["alg"])
		}
		publicKey := &m.secret.PublicKey
		return publicKey, nil
	})
	return token, err
}

// Валидацпия access токена
func (m *Middleware) IsTokenValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.ExtractToken(w, r)
		// в виду безопасности ошибка не уточняется
		if err != nil || !token.Valid {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		if token.Claims.(jwt.MapClaims)["type"].(string) != "access" {
			http.Error(w, "access not", http.StatusUnauthorized)
			return
		}
		_, ok := token.Claims.(jwt.MapClaims)["id"]
		if !ok {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		_, ok = token.Claims.(jwt.MapClaims)["is_org"]
		if !ok {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
