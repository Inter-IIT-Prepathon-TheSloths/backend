package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
)

func (uc *UserController) SignUp(c echo.Context) error {
	var userDetails models.LoginUserDetails

	if err := c.Bind(&userDetails); err != nil {
		return err
	}

	validation_err := utils.Validate(userDetails)
	if validation_err != nil {
		return validation_err
	}

	filter := utils.ConstructEmailFilter(userDetails.Email)
	existingUser, err := uc.service.GetUser(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Email already registered")
	}

	hashedPassword, err := utils.HashPassword(userDetails.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Emails:   []models.Email{{Email: userDetails.Email, IsVerified: false}},
		Password: hashedPassword,
	}

	id, err := uc.service.CreateUser(c.Request().Context(), user)
	if err != nil {
		return err
	}

	jwt, err := utils.CreateJwtToken(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to sign the jwt")
	}

	return c.JSON(http.StatusCreated, map[string]string{"jwt": jwt})
}

func (uc *UserController) Login(c echo.Context) error {
	var userDetails models.LoginUserDetails
	if err := c.Bind(&userDetails); err != nil {
		return err
	}

	filter := utils.ConstructEmailFilter(userDetails.Email)
	existingUser, err := uc.service.GetUser(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid Credentials")
	}

	if existingUser.Password == "" {
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?id=%s", fmt.Sprintf("%s/backend_redirect", os.Getenv("FRONTEND_URL")), existingUser.ID.Hex()))
	}

	err = utils.VerifyPassword(existingUser.Password, userDetails.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	jwt, err := utils.CreateJwtToken(existingUser.ID.Hex())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"jwt": jwt})
}
