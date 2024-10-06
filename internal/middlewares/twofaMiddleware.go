package middlewares

import (
	"net/http"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/controllers"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TwofaMiddleware(userController *controllers.UserController, sensitive bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user_id := c.Get("_id").(primitive.ObjectID)
			user := c.Get("user").(*models.User)
			twofa, err := userController.GetUserService().GetTwoFactor(c.Request().Context(), user_id)

			if !sensitive && !user.TwofaEnabled {
				c.Set("twofa", &models.TwoFactor{})
				return next(c)
			}

			// normal login
			if !sensitive && user.TwofaEnabled {
				normal_ok := c.Get("normal_ok").(bool)
				if normal_ok {
					c.Set("twofa", twofa)
					return next(c)
				}
			}

			if sensitive && user.TwofaEnabled {
				sensitive_ok := c.Get("sensitive_ok").(bool)
				if sensitive_ok {
					c.Set("sensitive_ok", true)
					return next(c)
				}
			}

			twofaCode := c.Request().Header.Get("X-2FA-Code")
			backupCode := false
			if twofaCode == "" {
				twofaCode = c.Request().Header.Get("X-2FA-Backup")
				if twofaCode == "" {
					return echo.NewHTTPError(http.StatusForbidden, "Please provide the totp code or backup code")
				}

				newBackupCodes, err := utils.UpdatedBackupCodes(twofaCode, twofa.BackupCodes)
				if err != nil {
					return err
				}

				twofa.BackupCodes = newBackupCodes
				err = userController.GetUserService().UpdateTwofa(c.Request().Context(), user_id, twofa)
				if err != nil {
					return err
				}

				backupCode = true
			}
			if twofaCode == "" {
				return echo.NewHTTPError(http.StatusForbidden, "Please provide the totp code")
			}

			if err != nil {
				return err
			}

			if !backupCode && !totp.Validate(twofaCode, twofa.Secret) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Two-factor authentication failed")
			}

			c.Set("twofa", twofa)
			return next(c)
		}
	}
}
