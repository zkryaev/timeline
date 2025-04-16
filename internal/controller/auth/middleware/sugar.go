package middleware

import (
	"context"
	"errors"
	"timeline/internal/controller/scope"
	"timeline/internal/entity"

	"github.com/golang-jwt/jwt/v5"
)

func GetTokenDataFromCtx(settings *scope.Settings, ctx context.Context) (entity.TokenData, error) {
	if !settings.EnableAuthorization {
		return entity.TokenData{}, nil
	}
	tdata, ok := ctx.Value(TokenData("token")).(entity.TokenData)
	if !ok {
		return tdata, errors.New("token not found in request context")
	}
	return tdata, nil
}

func GetTokenData(c jwt.Claims) entity.TokenData {
	data := entity.TokenData{}
	id := c.(jwt.MapClaims)["id"].(float64)
	data.IsOrg = c.(jwt.MapClaims)["is_org"].(bool)
	data.ID = int(id)
	return data
}
