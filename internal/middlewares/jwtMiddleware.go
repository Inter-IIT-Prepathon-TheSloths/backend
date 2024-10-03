package middlewares

import (
	"net/http"
	"strings"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/controllers"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthenticationMiddleware(userController *controllers.UserController) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
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

			user, err := userController.GetUserService().GetUser(c.Request().Context(), bson.M{"_id": oid})
			if err != nil {
				return err
			}
			if user == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
			}
			c.Set("user", user)
			c.Set("_id", oid)
			c.Set("twofa_ok", claims["twofa_ok"])

			return next(c)
		}
	}
}
