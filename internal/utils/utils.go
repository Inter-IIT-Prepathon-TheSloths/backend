package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/services"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckHealth(client *mongo.Client, t int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	defer cancel()

	err := client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
		return err
	}

	return nil
}

func GenerateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func Validate(s interface{}) error {
	validate := validator.New()
	err := validate.Struct(s)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Validation failed for field '%s', violating '%s' condition\n", err.Field(), err.Tag()))
		}
	}
	return nil
}

func GetEmailBody(email string, emails []models.Email) models.Email {
	for _, e := range emails {
		if e.Email == email {
			return e
		}
	}
	return models.Email{}
}

func ConstructEmailFilter(emails []models.Email) bson.M {
	var emailList []string

	for _, email := range emails {
		emailList = append(emailList, email.Email)
	}

	filter := bson.M{
		"emails.email": bson.M{
			"$in": emailList,
		},
	}
	return filter
}

func UpdateEmails(emailBody models.Email, emails []models.Email) []models.Email {
	for i, e := range emails {
		if e.Email == emailBody.Email {
			emails[i] = emailBody
			break
		}
	}
	return emails
}

func GenerateOTP() (string, error) {
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	// Generate a random number in the range [100000, 999999]
	otpBig, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return "", err
	}

	// Add the minimum value (100000) to get a 6-digit number
	otp := otpBig.Add(otpBig, min).Int64()

	// Return OTP as a string
	return fmt.Sprintf("%06d", otp), nil
}

func SendVerificationCode(code, email string) error {
	subject := "Email Verification Code - The Sloths"
	heading := "Verification Code"
	info1 := "To activate your email, please use the given OTP. Don't share with anyone :)"
	link := ""
	button_text := code
	time_duration := "15 Minutes"
	regenerate_link := os.Getenv("BACKEND_URL") + "/api/v1/auth/resend_code?email=" + email

	err := SendEmail([]string{email}, subject, heading, info1, link, button_text, time_duration, regenerate_link)
	return err
}

func SendRecoveryMail(link, email, time_duration string) error {
	subject := "Password Recovery Link - The Sloths"
	heading := "Password Recovery"
	info1 := "To reset your password, please click on the given link."
	button_text := "Recovery Link"
	regenerate_link := os.Getenv("FRONTEND_URL") + "/forgot_password" + email

	err := SendEmail([]string{email}, subject, heading, info1, link, button_text, time_duration, regenerate_link)
	return err
}

func VerifyRecoveryToken(token string, sv *services.UserService, c echo.Context) (*models.User, error) {
	token_parts, err := DecodeToken(token)
	if err != nil {
		return nil, err
	}
	id := token_parts[0]
	code := token_parts[1]
	fmt.Println(id)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user, err := sv.GetUser(c.Request().Context(), bson.M{"_id": oid})
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	if user.ForgotPassword.Code != code {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid verification code")
	}

	now := time.Now()
	if now.After(user.ForgotPassword.ExpiresAt) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Recovery code has expired")
	}

	return user, nil
}

func AtleastOneVerifiedEmailExists(emails []models.Email) bool {
	for _, e := range emails {
		if e.IsVerified {
			return true
		}
	}
	return false
}

func UpdatedBackupCodes(code string, backupCodes []string) ([]string, error) {
	if len(backupCodes) == 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Backup codes have been exhausted, sorry :(")
	}

	found := false
	for i, b := range backupCodes {
		if b == code {
			found = true
			backupCodes = append(backupCodes[:i], backupCodes[i+1:]...)
			break
		}
	}

	if !found {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid backup code")
	}
	return backupCodes, nil
}
