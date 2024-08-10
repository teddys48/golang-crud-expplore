package helper

import "github.com/golang-jwt/jwt/v5"

func GenerateToken(key interface{}, v jwt.Claims) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, v).SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}
