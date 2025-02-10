package validation

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrWrongClaims = errors.New("invalid token claims")
)

func IsCodeExpired(created_at time.Time) bool {
	currentTime := time.Now().UTC()
	if currentTime.Sub(created_at) > 5*time.Minute {
		return true
	}
	return false
}

// Проверяем что в формате UTC, что между текущей датой и датой создания не больше дня.
func IsAccountExpired(created_at time.Time) bool {
	Day := 24 * time.Hour
	currentDate := time.Now().UTC()                    // Получаем текущее время в UTC
	createdDate := created_at.Truncate(24 * time.Hour) // Обрезаем до суток
	if currentDate.Sub(createdDate) > Day {
		return true
	}
	return false
}

// Проверяем наличие полей и верного ли они типа
func ValidateTokenClaims(c jwt.Claims) error {
	if _, ok := c.(jwt.MapClaims)["id"].(float64); !ok {
		return ErrWrongClaims
	}
	if _, ok := c.(jwt.MapClaims)["is_org"].(bool); !ok {
		return ErrWrongClaims
	}
	if _, ok := c.(jwt.MapClaims)["type"].(string); !ok {
		return ErrWrongClaims
	}
	return nil
}
