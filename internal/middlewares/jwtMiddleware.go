package middlewares

import (
	"net/http"
	"strings"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearerString := c.Request().Header.Get("Authorization")
		if bearerString == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Please provide a bearer token")
		}
		words := strings.Split(bearerString, " ")
		if len(words) != 2 || words[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid authorization format")
		}

		token, err := utils.VerifyJwt(words[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid JWT Token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
		}

		oid, err := primitive.ObjectIDFromHex(claims["_id"].(string))
		if err != nil {
			return err
		}
		c.Set("_id", oid)

		return next(c)
	}
}
