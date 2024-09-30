package models

import "time"

type Verification struct {
	Email     string      `bson:"email" json:"email"`
	Code      string      `bson:"code" json:"code"`
	ExpiresAt time.Time   `bson:"expires_at"`
	Extras    interface{} `bson:"extras" json:"extras"`
}
