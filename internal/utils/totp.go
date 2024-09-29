package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/pquerna/otp/totp"
)

func Generate2faKey(id string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "TheSloths",
		AccountName: id,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Secret: %s\n", key.Secret())
	return key.Secret(), nil
}

func generateBackupCode() (string, error) {
	codeLength := 8
	bytes := make([]byte, codeLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateBackupCodes() ([]string, error) {
	var backupCodes []string
	for i := 0; i < 10; i++ {
		code, err := generateBackupCode()
		if err != nil {
			return nil, err
		}
		backupCodes = append(backupCodes, code)
	}
	return backupCodes, nil
}
