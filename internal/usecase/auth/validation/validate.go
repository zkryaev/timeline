package validation

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrWrongClaims    = errors.New("invalid token claims")
	ErrClaimsNotFound = errors.New("not found token claims")
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
func ValidateTokenClaims(req *jwt.Token) error {
	_, ok := req.Claims.(jwt.MapClaims)["id"].(float64)
	if !ok {
		log.Println("id!")
		return ErrWrongClaims
	}
	_, ok = req.Claims.(jwt.MapClaims)["is_org"].(bool)
	if !ok {
		log.Println("is_org!")
		return ErrWrongClaims
	}
	_, ok = req.Claims.(jwt.MapClaims)["type"].(string)
	if !ok {
		log.Println("type!")
		return ErrWrongClaims
	}
	return nil
}
