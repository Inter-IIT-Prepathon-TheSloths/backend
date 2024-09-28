package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwt_secret = []byte(os.Getenv("JWT_SECRET"))

func CreateJwtToken(_id string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id": _id,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // Expiry: 1 Week
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(jwt_secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJwt(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwt_secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing token issue")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
