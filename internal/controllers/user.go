package controllers

import (
	"net/http"

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

	jwt, err := utils.CreateJwtToken(user.ID.Hex())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"token": jwt})
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
