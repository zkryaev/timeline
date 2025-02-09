package secret

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"timeline/internal/libs/envars"
)

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	pathToSecret := envars.GetPathByEnv("SECRET_PATH")
	if pathToSecret == "" {
		return nil, fmt.Errorf("empty secret-path-env: %s", pathToSecret)
	}
	if _, err := os.Stat(pathToSecret); os.IsNotExist(err) {
		return nil, fmt.Errorf("file with secret does not exist: %w", err)
	}
	key, err := os.ReadFile(pathToSecret)
	if err != nil {
		fmt.Println(key)
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block from the key")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %v", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}

	return rsaPrivateKey, nil
}
