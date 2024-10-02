package models

import "time"

type Verification struct {
	Email     string            `bson:"email" json:"email"`
	Code      string            `bson:"code" json:"code"`
	ExpiresAt time.Time         `bson:"expires_at"`
	Extras    map[string]string `bson:"extras" json:"extras"`
}

type VerificationBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}
