package common

import (
	"timeline/internal/entity"

	"github.com/golang-jwt/jwt/v5"
)

func GetTokenData(c jwt.Claims) entity.TokenData {
	data := entity.TokenData{}
	id := c.(jwt.MapClaims)["id"].(float64)
	data.IsOrg = c.(jwt.MapClaims)["is_org"].(bool)
	data.ID = int(id)
	return data
}