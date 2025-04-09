package verification

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"time"

	"go.uber.org/zap"
)

var (
	minCode   int64 = 1000
	maxCode   int64 = 9999
	codeRange int64 = 1 + (maxCode - minCode)
)

func GenerateCode(logger *zap.Logger) (string, error) {
	max := big.NewInt(codeRange)
	if max.Sign() <= 0 {
		logger.Warn("big.NewInt returned number with sign <= 0", zap.Int64(max.Int64()), zap.Int64("codeRange", codeRange))
		max.SetInt64(minCode + time.Time.Unix())
	}
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Приводим результат к 4-значному числу и преобразуем его в строку
	code := int(n.Int64() + minCode)
	return strconv.Itoa(code), nil
}
