package helper

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func PrivateKey() (*rsa.PrivateKey, error) {
	getKey, err := os.ReadFile("private.key")
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(getKey)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func PublicKey() (*rsa.PublicKey, error) {
	getKey, err := os.ReadFile("public.key")
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(getKey)
	if err != nil {
		return nil, err
	}

	return key, nil
}
