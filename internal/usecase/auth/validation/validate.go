package validation

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrWrongClaims = errors.New("invalid token claims")
)

func IsCodeExpired(createdAt time.Time) bool {
	return time.Now().UTC().Sub(createdAt) > 5*time.Minute
}

// Проверяем что в формате UTC, что между текущей датой и датой создания не больше дня.
func IsAccountExpired(createdAt time.Time) bool {
	day := 24 * time.Hour
	currentDate := time.Now().UTC()        // Получаем текущее время в UTC
	createdDate := createdAt.Truncate(day) // Обрезаем до суток
	return currentDate.Sub(createdDate) > day
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
