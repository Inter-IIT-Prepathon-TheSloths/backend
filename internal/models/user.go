package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Emails       []Email            `bson:"emails" validate:"required,dive,required,dive,email" json:"emails"`
	Password     string             `bson:"password,omitempty" json:"-"`
	Picture      string             `bson:"picture,omitempty" json:"picture"`
	TwofaEnabled bool               `bson:"twofa_enabled"`
	CreatedAt    time.Time          `bson:"created_at" json:"-"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"-"`
}

type Email struct {
	Email      string `bson:"email" validate:"required,email" json:"email"`
	IsVerified bool   `bson:"is_verified" json:"is_verified"`
}

type LoginUserDetails struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type VerificationCode struct {
	Code      string    `bson:"code"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type GoogleUserDetails struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type UserPassword struct {
	ID       string `json:"id"`
	Password string `json:"password" validate:"required,min=6"`
}

type Signup struct {
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password"`
	Code      string    `bson:"code" json:"code"`
	ExpiresAt time.Time `bson:"expires_at"`
}
