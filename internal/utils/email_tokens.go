package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func EncodeToken(str string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	return encoded
}

func DecodeToken(token string) ([]string, error) {
	decodedBytes, err := base64.RawStdEncoding.DecodeString(token)
	if err != nil {
		return []string{}, err
	}
	fmt.Println("Token", token)

	decoded := string(decodedBytes)
	parts := strings.Split(decoded, ":")
	return parts, nil
}
