package controllers

import (
	"net/http"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *UserController) Generate2faSecret(c echo.Context) error {
	id := c.Get("_id").(primitive.ObjectID)
	user := c.Get("user").(*models.User)

	if user.TwofaEnabled {
		return echo.NewHTTPError(http.StatusUnauthorized, "TwoFactor Authentication already enabled")
	}

	key, err := utils.Generate2faKey(id.Hex())
	if err != nil {
		return err
	}

	twofa := &models.TwoFactor{
		UserID: id,
		Secret: key,
	}

	err = uc.service.UpdateTwofa(c.Request().Context(), id, twofa)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Secret created successfully",
		"secret":  key,
	})
}

func (uc *UserController) Regenerate2faSecret(c echo.Context) error {
	user := c.Get("user").(*models.User)
	if !user.TwofaEnabled {
		return echo.NewHTTPError(http.StatusBadRequest, "Twofa isn't enabled to be regenerated")
	}

	twofa := c.Get("twofa").(*models.TwoFactor)

	key, err := utils.Generate2faKey(twofa.UserID.Hex())
	if err != nil {
		return err
	}
	twofa.Secret = key

	err = uc.service.UpdateTwofa(c.Request().Context(), twofa.UserID, twofa)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Secret regenerated successfully",
		"secret":  key,
	})
}

func (uc *UserController) Disable2fa(c echo.Context) error {
	id := c.Get("_id").(primitive.ObjectID)
	user := c.Get("user").(*models.User)
	err := uc.service.DeleteTwoFactor(c.Request().Context(), id)
	if err != nil {
		return err
	}

	user.TwofaEnabled = false
	err = uc.service.UpdateUser(c.Request().Context(), id.Hex(), user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "2FA disabled successfully"})
}

func (uc *UserController) Enable2fa(c echo.Context) error {
	twofa := c.Get("twofa").(*models.TwoFactor)
	user := c.Get("user").(*models.User)

	if twofa.Secret == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide a valid totp")
	}

	backupCodes, err := utils.GenerateBackupCodes()
	if err != nil {
		return err
	}

	twofa.BackupCodes = backupCodes
	err = uc.service.UpdateTwofa(c.Request().Context(), twofa.UserID, twofa)
	if err != nil {
		return err
	}

	user.TwofaEnabled = true
	err = uc.service.UpdateUser(c.Request().Context(), user.ID.Hex(), user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "2FA enabled successfully"})
}

func (uc *UserController) GetTwofaInfo(c echo.Context) error {
	twofa := c.Get("twofa").(*models.TwoFactor)

	if twofa == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Two-factor authentication not found")
	}

	twofaToSend := struct {
		Secret      string   `json:"secret"`
		BackupCodes []string `json:"backup_codes"`
	}{
		Secret:      twofa.Secret,
		BackupCodes: twofa.BackupCodes,
	}

	return c.JSON(http.StatusOK, twofaToSend)
}

func (uc *UserController) RegenerateBackup2fa(c echo.Context) error {
	twofa := c.Get("twofa").(*models.TwoFactor)

	if twofa.Secret == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide a valid totp")
	}

	backupCodes, err := utils.GenerateBackupCodes()
	if err != nil {
		return err
	}
	twofa.BackupCodes = backupCodes

	err = uc.service.UpdateTwofa(c.Request().Context(), twofa.UserID, twofa)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Backup codes regenerated successfully",
		"backup_codes": backupCodes,
	})
}

func (uc *UserController) TwofaLogin(c echo.Context) error {
	id := c.Get("_id").(primitive.ObjectID)
	jwt, err := utils.CreateJwtToken(id.Hex(), true)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"token": jwt})
}
