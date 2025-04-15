package validation

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrWrongClaimsID   = errors.New("invalid token claims: id")
	ErrWrongClaimsOrg  = errors.New("invalid token claims: is_org")
	ErrWrongClaimsType = errors.New("invalid token claims: type")
	ErrNotAccessType   = errors.New("token type is not access")
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
		return fmt.Errorf("ValidateTokenClaims: %w", err)
	}
	if token.Claims.(jwt.MapClaims)["type"].(string) != "access" {
		return fmt.Errorf("%s: type=%q", ErrNotAccessType.Error(), token.Claims.(jwt.MapClaims)["type"].(string))
	}
	return nil
}

// Проверяем наличие полей и верного ли они типа
func ValidateTokenClaims(c jwt.Claims) error {
	if _, ok := c.(jwt.MapClaims)["id"].(float64); !ok {
		return ErrWrongClaimsID
	}
	if _, ok := c.(jwt.MapClaims)["is_org"].(bool); !ok {
		return ErrWrongClaimsOrg
	}
	if _, ok := c.(jwt.MapClaims)["type"].(string); !ok {
		return ErrWrongClaimsType
	}
	return nil
}
