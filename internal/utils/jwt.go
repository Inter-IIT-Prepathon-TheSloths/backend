package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/services"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwt_secret = []byte(config.JwtSecret)

func CreateJwtToken(_id string, twofa_ok, sensitive bool) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id":       _id,
		"twofa_ok":  twofa_ok,
		"sensitive": sensitive,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
		"iat":       time.Now().Unix(),
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

func HandleJwt(bearerString string, sensitive bool, sv *services.UserService, c context.Context) (*models.User, jwt.MapClaims, error) {
	if bearerString == "" {
		return nil, nil, echo.NewHTTPError(http.StatusBadRequest, "Please provide a bearer token")
	}
	words := strings.Split(bearerString, " ")
	if len(words) != 2 || words[0] != "Bearer" {
		return nil, nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid authorization format")
	}

	token, err := VerifyJwt(words[1])
	if err != nil {
		return nil, nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid JWT Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
	}

	if claims["sensitive"] != sensitive {
		return nil, nil, echo.NewHTTPError(http.StatusUnauthorized, "Wrong jwt token provided for login")
	}

	oid, err := primitive.ObjectIDFromHex(claims["_id"].(string))
	if err != nil {
		return nil, nil, err
	}

	user, err := sv.GetUser(c, bson.M{"_id": oid})
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, echo.NewHTTPError(http.StatusUnauthorized, "User not found")
	}

	return user, claims, nil
}
