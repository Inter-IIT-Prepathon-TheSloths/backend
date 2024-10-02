package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
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
