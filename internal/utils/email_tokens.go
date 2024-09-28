package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/google/uuid"
)

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

func GenerateVerificationCode() models.VerificationCode {
	token := uuid.New().String()
	expiration := time.Now().Add(24 * time.Hour)
	code := models.VerificationCode{
		Code:      token,
		ExpiresAt: expiration,
	}
	return code
}
