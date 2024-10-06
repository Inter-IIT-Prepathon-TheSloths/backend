package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
)

func (uc *UserController) SendVerificationMail(c echo.Context) error {
	var body models.VerificationBody
	if err := c.Bind(&body); err != nil {
		return err
	}
	validation_err := utils.Validate(body)
	if validation_err != nil {
		return validation_err
	}

	emailFilter := utils.ConstructEmailFilter([]models.Email{
		{Email: body.Email},
	})
	user, err := uc.service.GetUser(c.Request().Context(), emailFilter)
	if err != nil {
		return err
	}

	use := c.Param("use")

	var subject, heading, info1, time_duration, regenerate_link string
	var expiry_duration time.Duration
	var extras map[string]string
	expiry_duration = time.Hour

	if use == "signup" {
		if user != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Email already exists")
		}

		if len(body.Password) < 6 {
			return echo.NewHTTPError(http.StatusBadRequest, "Password should be at least 6 characters long")
		}

		subject = "Email verification OTP"
		heading = "Welcome to Bizsights!"
		info1 = "To complete your registration, please use this OTP given below"
		time_duration = "1 Hour"
		regenerate_link = fmt.Sprintf("%s/signup", config.FrontendUrl)
		hashedPassword, err := utils.HashPassword(body.Password)
		if err != nil {
			return err
		}
		extras = map[string]string{
			"Password": hashedPassword,
		}
	} else if use == "reset_password" {
		if user == nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		subject = "Password reset OTP"
		heading = "Reset your password"
		info1 = "To reset your password, please use this OTP given below"
		time_duration = "1 Hour"
		regenerate_link = fmt.Sprintf("%s/reset_password", config.FrontendUrl)
		extras = nil
	} else if use == "verify_email" {
		if user == nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		emailBody := utils.GetEmailBody(body.Email, user.Emails)
		if emailBody.IsVerified {
			return echo.NewHTTPError(http.StatusNotFound, "Email already verified")
		}
		subject = "Email verification OTP"
		heading = "Verify your email address"
		info1 = "To verify your email address, please use this OTP given below"
		time_duration = "1 Hour"
		regenerate_link = fmt.Sprintf("%s/emails", config.FrontendUrl)
		extras = nil
	} else if use == "update_password" {
		if user == nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		if len(body.Password) < 6 {
			return echo.NewHTTPError(http.StatusBadRequest, "Password should be at least 6 characters long")
		}

		subject = "Password Update OTP"
		heading = "Update your password"
		info1 = "To update your password, please use this OTP given below"
		time_duration = "1 Hour"
		regenerate_link = fmt.Sprintf("%s/update_password", config.FrontendUrl)
		hashedPassword, err := utils.HashPassword(body.Password)
		if err != nil {
			return err
		}
		extras = map[string]string{
			"Password": hashedPassword,
		}
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid use parameter")
	}

	if err := utils.SendVerification(c, uc.service, body.Email, subject, heading, info1, time_duration, regenerate_link, expiry_duration, extras); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Verification email sent successfully"})
}

func (uc *UserController) VerifyVerificationCode(c echo.Context) error {
	var body struct {
		Email    string `json:"email" validate:"email"`
		Code     string `json:"code"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	validation_err := utils.Validate(body)
	if validation_err != nil {
		return validation_err
	}

	verif, err := uc.service.GetVerification(c.Request().Context(), body.Email)
	if err != nil {
		return err
	}

	if verif.Code != body.Code || time.Now().After(verif.ExpiresAt) {
		return echo.NewHTTPError(http.StatusUnauthorized, "code is invalid or expired")
	}

	var errToSend error

	use := c.Param("use")
	if use == "signup" {
		user := &models.User{
			Emails: []models.Email{
				{Email: body.Email, IsVerified: true},
			},
			Password: verif.Extras["Password"],
		}

		_, err = uc.service.CreateUser(c.Request().Context(), user)
		if err != nil {
			return err
		}

		errToSend = c.JSON(http.StatusCreated, map[string]string{
			"message": "Signed up successfully",
		})

	} else {
		emailFilter := utils.ConstructEmailFilter([]models.Email{
			{Email: body.Email, IsVerified: true},
		})
		user, err := uc.service.GetUser(c.Request().Context(), emailFilter)
		if err != nil {
			return err
		}

		if user == nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		emailBody := utils.GetEmailBody(body.Email, user.Emails)
		if emailBody.Email == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Email not found")
		}

		if !emailBody.IsVerified {
			emailBody.IsVerified = true
			user.Emails = utils.UpdateEmails(emailBody, user.Emails)
		}

		if use == "reset_password" {

			if len(body.Password) < 6 {
				return echo.NewHTTPError(http.StatusBadRequest, "Password should be at least 6 characters long")
			}

			hashedPassword, err := utils.HashPassword(body.Password)
			if err != nil {
				return err
			}
			user.Password = hashedPassword

			errToSend = c.JSON(http.StatusOK, map[string]string{
				"message": "Password reset successfully",
			})
		} else if use == "verify_email" {
			emailBody := utils.GetEmailBody(body.Email, user.Emails)
			if emailBody.Email == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Email not found")
			}

			if emailBody.IsVerified {
				return echo.NewHTTPError(http.StatusNotFound, "Email already verified")
			}

			errToSend = c.JSON(http.StatusOK, map[string]string{
				"message": "Email verification successful",
			})
		} else if use == "update_password" {
			user.Password = verif.Extras["Password"]

			errToSend = c.JSON(http.StatusOK, map[string]string{
				"message": "Password updated successfully",
			})
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid use parameter")
		}

		if err = uc.service.UpdateUser(c.Request().Context(), user.ID.Hex(), user); err != nil {
			return err
		}
	}

	if err = uc.service.DeleteVerification(c.Request().Context(), body.Email); err != nil {
		return err
	}

	return errToSend
}
