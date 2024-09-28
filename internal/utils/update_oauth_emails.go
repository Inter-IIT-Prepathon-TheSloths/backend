package utils

import "github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"

func UpdateOauthEmails(existingEmails, newEmails []models.Email) ([]models.Email, bool) {
	change := false

	// Create a map to check for existing emails quickly
	emailMap := make(map[string]bool)
	for _, email := range existingEmails {
		emailMap[email.Email] = email.IsVerified
	}

	// Iterate through new emails
	for _, newEmail := range newEmails {
		if verified, exists := emailMap[newEmail.Email]; exists {
			// If the new email exists and is not verified, mark it as verified
			if !verified {
				// Find and update the existing email
				for i, existingEmail := range existingEmails {
					if existingEmail.Email == newEmail.Email {
						existingEmails[i].IsVerified = true
						change = true
					}
				}
			}
		} else {
			// If the new email doesn't exist, append it
			existingEmails = append(existingEmails, newEmail)
			change = true
		}
	}

	return existingEmails, change
}
