package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
)

func (uc *UserController) SignUp(c echo.Context) error {
	var userDetails models.LoginUserDetails

	// Reading user details from request body
	if err := c.Bind(&userDetails); err != nil {
		return err
	}

	// Validating email and password
	validation_err := utils.Validate(userDetails)
	if validation_err != nil {
		return validation_err
	}

	// Get the user with the given email
	filter := utils.ConstructEmailFilter([]models.Email{
		{Email: userDetails.Email, IsVerified: true},
	})
	existingUser, err := uc.service.GetUser(c.Request().Context(), filter)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Email already registered")
	}

	// hash the password
	hashedPassword, err := utils.HashPassword(userDetails.Password)
	if err != nil {
		return err
	}

	// Create a signup verification document in DB
	code, err := utils.GenerateOTP()
	if err != nil {
		return err
	}
	signup := &models.Signup{
		Email:     userDetails.Email,
		Code:      code,
		Password:  hashedPassword,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	err = uc.service.UpdateSignup(c.Request().Context(), userDetails.Email, signup)
	if err != nil {
		return err
	}

	// Send verification code on email
	err = utils.SendVerificationCode(code, userDetails.Email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Verification code sent on email"})
}

func (uc *UserController) Login(c echo.Context) error {
	var userDetails models.LoginUserDetails
	if err := c.Bind(&userDetails); err != nil {
		return err
	}

	filter := utils.ConstructEmailFilter([]models.Email{
		{Email: userDetails.Email, IsVerified: true},
	})
	existingUser, err := uc.service.GetUser(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid Credentials")
	}

	// If the user logged in with Oauth but didn't create password, so redirect him to first create password
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
