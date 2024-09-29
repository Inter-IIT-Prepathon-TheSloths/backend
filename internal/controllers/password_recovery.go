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

func (uc *UserController) ForgotPassword(c echo.Context) error {
	email := c.QueryParam("email")

	emailFilter := utils.ConstructEmailFilter([]models.Email{{Email: email}})
	user, err := uc.service.GetUser(c.Request().Context(), emailFilter)
	if err != nil {
		return err
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	code := utils.GenerateVerificationCode(time.Hour)

	user.ForgotPassword = code
	err = uc.service.UpdateUser(c.Request().Context(), user.ID.Hex(), user)
	if err != nil {
		return err
	}

	token_str := fmt.Sprintf("%s:%s", user.ID.Hex(), code.Code)
	encodedToken := utils.EncodeToken(token_str)
	verificationLink := fmt.Sprintf("%s/forgot_password?token=%s", os.Getenv("FRONTEND_URL"), encodedToken)
	err = utils.SendRecoveryMail(verificationLink, email, "1 Hour")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Verification link has been sent"})
}

func (uc *UserController) VerifyForgotPassword(c echo.Context) error {
	token := c.Param("token")
	_, err := utils.VerifyRecoveryToken(token, uc.service, c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password recovery code verified"})
}

func (uc *UserController) RecoverPassword(c echo.Context) error {
	var changePassword struct {
		Token    string `json:"token"`
		Password string `json:"password" validate:"min=6"`
	}

	if err := c.Bind(&changePassword); err != nil {
		return err
	}

	user, err := utils.VerifyRecoveryToken(changePassword.Token, uc.service, c)
	if err != nil {
		return err
	}

	err = utils.Validate(changePassword)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(changePassword.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.ForgotPassword = models.VerificationCode{}
	err = uc.service.UpdateUser(c.Request().Context(), user.ID.Hex(), user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password changed successfully"})
}
