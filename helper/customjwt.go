package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAdminJWT(signKey string, adminID string) (string, error) {
	claims := jwt.MapClaims{
		"id":    adminID,
		"admin": true,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(signKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
