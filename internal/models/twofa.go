package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TwoFactor struct {
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Secret      string             `bson:"secret" json:"secret"`
	BackupCodes []string           `bson:"backup_codes" json:"backup_codes"`
}
