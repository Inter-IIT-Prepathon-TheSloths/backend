package utils

import (
	"fmt"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/golang-jwt/jwt"
)

var jwt_secret = []byte(config.JwtSecret)

func CreateJwtToken(_id string, twofa_ok bool) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":      _id,
		"twofa_ok": twofa_ok,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
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
