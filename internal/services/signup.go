package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *UserService) getSignupsCollection() *mongo.Collection {
	return s.client.Database(os.Getenv("DB_NAME")).Collection("signups")
}

func (s *UserService) getVerificationsCollection() *mongo.Collection {
	return s.client.Database(os.Getenv("DB_NAME")).Collection("signups")
}

func (s *UserService) CreateSignup(ctx context.Context, email, code, password string) error {
	signup := &models.Signup{
		Email:     email,
		Code:      code,
		Password:  password,
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}

	_, err := s.getSignupsCollection().InsertOne(ctx, signup)
	return err
}

func (s *UserService) GetSignup(ctx context.Context, email string) (*models.Signup, error) {
	var signup models.Signup
	err := s.getSignupsCollection().FindOne(ctx, bson.M{"email": email}).Decode(&signup)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no signup found")
	}
	if err != nil {
		return nil, fmt.Errorf("error Decoding User Doc")
	}

	return &signup, nil
}

func (s *UserService) UpdateSignup(ctx context.Context, email string, signup *models.Signup) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": signup}

	opts := options.Update().SetUpsert(true)
	_, err := s.getSignupsCollection().UpdateOne(ctx, filter, update, opts)
	return err
}

func (s *UserService) DeleteSignup(ctx context.Context, email string) error {
	filter := bson.M{"email": email}
	_, err := s.getSignupsCollection().DeleteOne(ctx, filter)
	return err
}
