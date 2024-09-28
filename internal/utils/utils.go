package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var jwt_secret = []byte(os.Getenv("JWT_SECRET"))

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

func CreateJwtToken(_id string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id": _id,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // Expiry: 1 Week
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(jwt_secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwt_secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing token issue")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

func GenerateState() (string, error) {
	b := make([]byte, 16) // 16 bytes == 128 bits
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

func GenerateVerificationCode() models.VerificationCode {
	token := uuid.New().String()
	expiration := time.Now().Add(24 * time.Hour)
	code := models.VerificationCode{
		Code:      token,
		ExpiresAt: expiration,
	}
	return code
}

func EncodeToken(id, email, code string) string {
	combined := fmt.Sprintf("%s:%s:%s", id, email, code)

	encoded := base64.StdEncoding.EncodeToString([]byte(combined))
	return encoded
}

func DecodeToken(token string) (string, string, string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", "", "", err
	}

	decoded := string(decodedBytes)
	parts := strings.Split(decoded, ":")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid token format")
	}
	return parts[0], parts[1], parts[2], nil
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
