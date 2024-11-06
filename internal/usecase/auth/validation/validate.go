package validation

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrWrongClaims    = errors.New("invalid token claims")
	ErrClaimsNotFound = errors.New("not found token claims")
)

// пароль там должен быть не меньше 8 символов, состоять там из таких то таких символов и не должен состоять из других символов
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
