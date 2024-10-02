package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *UserService) getVerificationsCollection() *mongo.Collection {
	return s.client.Database(config.DbName).Collection("verifications")
}

func (s *UserService) GetVerification(ctx context.Context, email string) (*models.Verification, error) {
	var verif models.Verification
	err := s.getVerificationsCollection().FindOne(ctx, bson.M{"email": email}).Decode(&verif)
	if err == mongo.ErrNoDocuments {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "code is invalid or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("error Decoding User Doc")
	}

	return &verif, nil
}

func (s *UserService) UpdateVerification(ctx context.Context, email, code string, expires_at time.Time, extras map[string]string) error {
	filter := bson.M{"email": email}

	verif := &models.Verification{
		Email:     email,
		Code:      code,
		ExpiresAt: expires_at,
		Extras:    extras,
	}
	update := bson.M{"$set": verif}

	opts := options.Update().SetUpsert(true)
	_, err := s.getVerificationsCollection().UpdateOne(ctx, filter, update, opts)
	return err
}

func (s *UserService) DeleteVerification(ctx context.Context, email string) error {
	filter := bson.M{"email": email}
	_, err := s.getVerificationsCollection().DeleteOne(ctx, filter)
	return err
}
