package secret

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadPrivateKey(env string) (*rsa.PrivateKey, error) {
	key := os.Getenv(env)
	if key == "" {
		return nil, fmt.Errorf("environment variable %s is not set, got empty key %s", env, key)
	}
	// TODO: Как ключ генерить то?
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block from the key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %v", err)
	}

	return privateKey, nil
}
