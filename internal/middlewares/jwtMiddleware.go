package middlewares

import (
	"net/http"
	"strings"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
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

		token, err := utils.VerifyToken(words[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid JWT Token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
		}

		c.Set("_id", claims["_id"])

		return next(c)
	}
}
