package validation

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrWrongClaims    = errors.New("invalid token claims")
	ErrClaimsNotFound = errors.New("not found token claims")
)

// Проверяем что дата создания <= текущая дата
func IsAccountExpired(created_at time.Time) bool {
	// обнуляет время, оставляя только дату
	currentDate := time.Now().Truncate(24 * time.Hour)
	createdDate := created_at.Truncate(24 * time.Hour)
	if currentDate.After(createdDate) {
		return true
	}
	return false
}

// Проверяем наличие полей и верного ли они типа
func ValidateTokenClaims(req *jwt.Token) error {
	_, ok := req.Claims.(jwt.MapClaims)["id"].(uint64)
	if !ok {
		return ErrWrongClaims
	}
	_, ok = req.Claims.(jwt.MapClaims)["is_org"].(bool)
	if !ok {
		return ErrWrongClaims
	}
	_, ok = req.Claims.(jwt.MapClaims)["type"].(string)
	if !ok {
		return ErrWrongClaims
	}
	return nil
}
