package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	UserId    primitive.ObjectID `bson:"user_id"`
	Token     string             `bson:"token"`
	ExpiresAt time.Time          `bson:"expires_at"`
}
