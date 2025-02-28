package verification

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

func GenerateCode() (string, error) {
	min, max := 1000, 9999 // nolint:revive -- faster
	rangeVal := max - min + 1

	n, err := rand.Int(rand.Reader, big.NewInt(int64(rangeVal)))
	if err != nil {
		return "", err
	}

	// Приводим результат к 4-значному числу и преобразуем его в строку
	code := int(n.Int64()) + min
	return strconv.Itoa(code), nil
}
