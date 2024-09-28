package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
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

func ConstructEmailFilter(email string) bson.M {
	filter := bson.M{
		"emails": bson.M{
			"$elemMatch": bson.M{
				"email": email,
			},
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
