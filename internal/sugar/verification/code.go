package verification

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

var (
	minCode   int64 = 1000
	maxCode   int64 = 9999
	codeRange int64 = 1 + (maxCode - minCode)
)

func GenerateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(codeRange))
	if err != nil {
		return "", err
	}

	// Приводим результат к 4-значному числу и преобразуем его в строку
	code := int(n.Int64() + minCode)
	return strconv.Itoa(code), nil
}
