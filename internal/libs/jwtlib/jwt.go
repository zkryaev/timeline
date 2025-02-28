package jwtlib

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"time"
	"timeline/internal/config"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/authdto"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidTokenType  = errors.New("invalid token type")
	ErrInvalidSignMethod = errors.New("invalid sign method")
)

func NewTokenPair(secret *rsa.PrivateKey, cfg config.Token, metadata *entity.TokenMetadata) (*authdto.TokenPair, error) {
	accessToken, err := NewToken(secret, cfg, metadata, "access")
	if err != nil {
		return nil, err
	}
	refreshToken, err := NewToken(secret, cfg, metadata, "refresh")
	if err != nil {
		return nil, err
	}
	return &authdto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func NewToken(secret interface{}, cfg config.Token, metadata *entity.TokenMetadata, tokenType string) (string, error) {
	var exp int64
	switch tokenType {
	case "access":
		exp = time.Now().Add(cfg.AccessTTL).Unix()
	case "refresh":
		exp = time.Now().Add(cfg.RefreshTTL).Unix()
	default:
		return "", ErrInvalidTokenType
	}
	var signMethod jwt.SigningMethod
	switch secret.(type) {
	case *rsa.PrivateKey:
		signMethod = jwt.SigningMethodRS256
	case *ecdsa.PrivateKey:
		signMethod = jwt.SigningMethodES256
	default:
		return "", ErrInvalidSignMethod
	}
	token := jwt.NewWithClaims(signMethod, jwt.MapClaims{
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
