package utils

import (
	Authentication "KAIA-BE/controllers/auth"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func ParseJWT(tokenstr string) (*Authentication.JWTClaim, error) {
	claims := &Authentication.JWTClaim{}
	token, err := jwt.ParseWithClaims(tokenstr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
