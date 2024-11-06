package middleware

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	secret *rsa.PrivateKey
}

func New(key *rsa.PrivateKey) *Middleware {
	return &Middleware{
		secret: key,
	}
}

var (
	ErrTokenNotFound   = errors.New("token not found")
	ErrAuthHeaderEmpty = errors.New("auth header empty")
)

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
		return m.secret, nil
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
		if token.Claims.(jwt.MapClaims)["type"] != "access" {
			http.Error(w, "", http.StatusUnauthorized)
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

// Валидация refresh токена
func (m *Middleware) IsRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.ExtractToken(w, r)
		if err != nil || !token.Valid {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		if token.Claims.(jwt.MapClaims)["type"] != "refresh" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
