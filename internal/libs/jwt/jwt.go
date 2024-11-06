package jwtlib

import (
	"crypto/rsa"
	"errors"
	"time"
	"timeline/internal/config"
	"timeline/internal/model"
	"timeline/internal/model/dto"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidTokenType = errors.New("invalid token type")
)

func NewTokenPair(secret *rsa.PrivateKey, cfg config.Token, metadata *model.TokenMetadata) (*dto.TokenPair, error) {
	AccessToken, err := NewToken(secret, cfg, metadata, "access")
	if err != nil {
		return nil, err
	}
	RefreshToken, err := NewToken(secret, cfg, metadata, "refresh")
	if err != nil {
		return nil, err
	}
	return &dto.TokenPair{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
	}, nil
}

func NewToken(secret *rsa.PrivateKey, cfg config.Token, metadata *model.TokenMetadata, tokenType string) (string, error) {
	var exp int64
	switch tokenType {
	case "access":
		exp = time.Now().Add(cfg.AccessTTL).Unix()
	case "refresh":
		exp = time.Now().Add(cfg.RefreshTTL).Unix()
	default:
		return "", ErrInvalidTokenType
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":     metadata.ID,
		"is_org": metadata.IsOrg,
		"type":   tokenType,
		"exp":    exp,
	})
	tokenEncoded, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenEncoded, nil
}
