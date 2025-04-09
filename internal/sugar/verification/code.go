package verification

import (
	"crypto/rand"
	"math/big"
	"strconv"

	"go.uber.org/zap"
)

var (
	minCode   int64 = 1000
	maxCode   int64 = 9999
	codeRange int64 = 1 + (maxCode - minCode)
)

func GenerateCode(logger *zap.Logger) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("GenerateCode", zap.Any("recover", r), zap.Int64("codeRange", codeRange))
		}
	}()
	n, err := rand.Int(rand.Reader, big.NewInt(codeRange))
	if err != nil {
		return "", err
	}

	// Приводим результат к 4-значному числу и преобразуем его в строку
	code := int(n.Int64() + minCode)
	return strconv.Itoa(code), nil
}
