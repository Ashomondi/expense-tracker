package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"expense-tracker/models"
)

var JwtKey = []byte("supersecretkey")

func ValidateToken(r *http.Request) (string, error) {
	tokenStr := r.Header.Get("Authorization")

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&models.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		},
	)

	if err != nil {
		return "", err
	}

	claims := token.Claims.(*models.Claims)

	return claims.Username, nil
}