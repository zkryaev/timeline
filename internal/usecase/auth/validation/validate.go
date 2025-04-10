package validation

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrWrongClaims   = errors.New("invalid token claims")
	ErrNotAccessType = errors.New("token type is not \"access\"")
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

func IsTokenValid(token *jwt.Token) (err error) {
	if !token.Valid {
		return fmt.Errorf("token.Valid : %v", token.Valid)
	}
	if err = ValidateTokenClaims(token.Claims); err != nil {
		return fmt.Errorf("%s: %w", ErrWrongClaims.Error(), err)
	}
	if token.Claims.(jwt.MapClaims)["type"].(string) != "access" {
		return fmt.Errorf("%s: type=%q", ErrNotAccessType.Error(), token.Claims.(jwt.MapClaims)["type"].(string))
	}
	return nil
}

// Проверяем наличие полей и верного ли они типа
func ValidateTokenClaims(c jwt.Claims) error {
	if _, ok := c.(jwt.MapClaims)["id"].(float64); !ok {
		return fmt.Errorf("\"id\" invalid")
	}
	if _, ok := c.(jwt.MapClaims)["is_org"].(bool); !ok {
		return fmt.Errorf("\"is_org\" invalid")
	}
	if _, ok := c.(jwt.MapClaims)["type"].(string); !ok {
		return fmt.Errorf("\"type\" invalid")
	}
	return nil
}
