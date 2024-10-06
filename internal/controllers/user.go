package controllers

import (
	"net/http"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc *UserController) CreatePassword(c echo.Context) error {
	var password models.UserPassword
	if err := c.Bind(&password); err != nil {
		return err
	}

	validation_error := utils.Validate(password)
	if validation_error != nil {
		return validation_error
	}

	oid, err := primitive.ObjectIDFromHex(password.ID)
	if err != nil {
		return err
	}

	user, err := uc.service.GetUser(c.Request().Context(), bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	hashedPassword, err := utils.HashPassword(password.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	err = uc.service.UpdateUser(c.Request().Context(), password.ID, user)
	if err != nil {
		return err
	}

	jwt, err := utils.CreateJwtToken(user.ID.Hex(), false, false, 15*time.Minute)
	if err != nil {
		return err
	}

	refreshToken, err := utils.CreateSessionToken(user.ID, false, false, 15*24*time.Hour, c.Request().Context(), uc.service)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"message":      "Logged In",
		"token":        jwt,
		"refreshToken": refreshToken,
		"status":       "success",
	})

	// return c.JSON(http.StatusCreated, map[string]string{"token": jwt})
}

func (uc *UserController) GetMyDetails(c echo.Context) error {
	user := c.Get("user").(*models.User)
	return c.JSON(http.StatusOK, user)
}

func (uc *UserController) AddEmail(c echo.Context) error {
	var email models.Email
	if err := c.Bind(&email); err != nil {
		return err
	}

	validation_error := utils.Validate(email)
	if validation_error != nil {
		return validation_error
	}

	user := c.Get("user").(*models.User)

	emailExisting := utils.GetEmailBody(email.Email, user.Emails)
	if emailExisting.Email != "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email already exists")
	}

	if len(user.Emails) >= 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "Maximum number of emails reached")
	}

	email.IsVerified = false
	user.Emails = append(user.Emails, email)

	err := uc.service.UpdateUser(c.Request().Context(), user.ID.Hex(), user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Email added successfully"})
}

func (uc *UserController) RefreshToken(c echo.Context) error {
	bearerString := c.Request().Header.Get("X-Refresh-Token")
	if bearerString == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token is required")
	}

	user, claims, refreshToken, err := utils.HandleJwt(bearerString, "Refresh", false, uc.service, c.Request().Context())
	if err != nil {
		return err
	}

	_, err = uc.service.GetSession(c.Request().Context(), user.ID, refreshToken)
	if err != nil {
		return err
	}

	twofa_ok := claims["twofa_ok"].(bool)

	jwt, err := utils.CreateJwtToken(user.ID.Hex(), twofa_ok, false, 15*time.Minute)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Token refreshed successfully",
		"token":   jwt,
	})
}

func (uc *UserController) Logout(c echo.Context) error {
	id := c.Get("_id").(primitive.ObjectID)
	refreshToken := c.Request().Header.Get("X-Refresh-Token")
	if refreshToken == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token is required")
	}

	if err := uc.service.DeleteSession(c.Request().Context(), id, refreshToken); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (uc *UserController) LogoutAll(c echo.Context) error {
	id := c.Get("_id").(primitive.ObjectID)
	if err := uc.service.DeleteAllSessions(c.Request().Context(), id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out from all devices successfully"})
}
