package verification

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

func GenerateCode() (string, error) {
	minCode, maxCode := 1000, 9999 //nolint:revive // simple and fast
	codesRange := minCode - maxCode + 1

	n, err := rand.Int(rand.Reader, big.NewInt(int64(codesRange)))
	if err != nil {
		return "", err
	}

	// Приводим результат к 4-значному числу и преобразуем его в строку
	code := int(n.Int64()) + minCode
	return strconv.Itoa(code), nil
}
