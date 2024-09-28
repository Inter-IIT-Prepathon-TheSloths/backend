package controllers

import (
	"fmt"
	"net/http"
	"os"
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

	jwt, err := utils.CreateJwtToken(user.ID.Hex())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"token": jwt})
}

func (uc *UserController) GetMyDetails(c echo.Context) error {
	id := c.Get("_id")
	user, err := uc.service.GetUser(c.Request().Context(), bson.M{"_id": id})
	if err != nil {
		return err
	}
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

	id := c.Get("_id").(primitive.ObjectID)

	user, err := uc.service.GetUser(c.Request().Context(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	emailExisting := utils.GetEmailBody(email.Email, user.Emails)
	if emailExisting.Email != "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email already exists")
	}

	if len(user.Emails) >= 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "Maximum number of emails reached")
	}

	email.IsVerified = false
	user.Emails = append(user.Emails, email)

	err = uc.service.UpdateUser(c.Request().Context(), id.Hex(), user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Email added successfully"})
}

func (uc *UserController) SendActivationMail(c echo.Context) error {
	email := c.Param("email")
	id := c.Get("_id").(primitive.ObjectID)

	user, err := uc.service.GetUser(c.Request().Context(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	emailBody := utils.GetEmailBody(email, user.Emails)
	if emailBody.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email not found")
	}

	if emailBody.IsVerified {
		return echo.NewHTTPError(http.StatusBadRequest, "Email already verified")
	}

	newCode := utils.GenerateVerificationCode()
	for ind, e := range user.Emails {
		if e.Email == email {
			user.Emails[ind].VerificationCode = newCode
			break
		}
	}

	err = uc.service.UpdateUser(c.Request().Context(), id.Hex(), user)
	if err != nil {
		return err
	}

	baseUrl := os.Getenv("FRONTEND_URL")
	subject := "Activate your account - The Sloths"
	heading := "Activate your account"
	info1 := "To activate your account, please click the button below and follow the instructions provided."
	link := fmt.Sprintf("%s/activate_account?token=%s", baseUrl, utils.EncodeToken(id.Hex(), email, newCode.Code))
	time_duration := "1 day"
	regenerate_link := os.Getenv("BACKEND_URL") + "/api/v1/auth/send_activation"

	err = utils.SendEmail([]string{email}, subject, heading, info1, link, time_duration, regenerate_link)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Activation link has been sent successfully to your email"})
}

func (uc *UserController) VerifyEmail(c echo.Context) error {
	token := c.Param("token")
	id, email, code, err := utils.DecodeToken(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid token")
	}

	oid, err := primitive.ObjectIDFromHex(id)
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

	emailBody := utils.GetEmailBody(email, user.Emails)
	if emailBody.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email not found")
	}

	if emailBody.VerificationCode.Code != code {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid verification code")
	}

	now := time.Now()
	if now.After(emailBody.VerificationCode.ExpiresAt) {
		return echo.NewHTTPError(http.StatusBadRequest, "Verification code has expired")
	}

	emailBody.IsVerified = true
	emailBody.VerificationCode = models.VerificationCode{}

	updatedEmails := utils.UpdateEmails(emailBody, user.Emails)
	user.Emails = updatedEmails

	err = uc.service.UpdateUser(c.Request().Context(), id, user)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Email verified successfully"})
}
