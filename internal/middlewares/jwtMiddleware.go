package middlewares

import (
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/controllers"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
)

func AuthenticationMiddleware(userController *controllers.UserController, sensitive bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var bearerString string
			if sensitive {
				bearerString = c.Request().Header.Get("X-2FA-Sensitive-Token")
			} else {
				bearerString = c.Request().Header.Get("Authorization")
			}

			if sensitive && bearerString == "" {
				c.Set("sensitive_ok", false)
				return next(c)
			}

			user, claims, _, err := utils.HandleJwt(bearerString, "JWT", sensitive, userController.GetUserService(), c.Request().Context())
			if err != nil {
				return err
			}

			c.Set("user", user)
			c.Set("_id", user.ID)
			var key string
			if !sensitive {
				key = "normal_ok"
			} else {
				key = "sensitive_ok"
			}
			c.Set(key, claims["twofa_ok"])
			return next(c)
		}
	}
}
