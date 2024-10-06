package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *UserService) getSessionsCollection() *mongo.Collection {
	return s.client.Database(config.DbName).Collection("sessions")
}

func (s *UserService) GetSession(ctx context.Context, userId primitive.ObjectID, refreshToken string) (*models.Session, error) {
	filter := bson.M{"user_id": userId, "token": refreshToken}

	var session models.Session
	err := s.getSessionsCollection().FindOne(ctx, filter).Decode(&session)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no sessions found")
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *UserService) CreateSession(ctx context.Context, userId primitive.ObjectID, refreshToken string, dur time.Duration) error {
	session := &models.Session{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(dur),
	}
	_, err := s.getSessionsCollection().InsertOne(ctx, session)
	return err
}

func (s *UserService) DeleteSession(ctx context.Context, userId primitive.ObjectID, refreshToken string) error {
	filter := bson.M{"user_id": userId, "token": refreshToken}
	_, err := s.getSessionsCollection().DeleteOne(ctx, filter)
	return err
}

func (s *UserService) DeleteAllSessions(ctx context.Context, userId primitive.ObjectID) error {
	filter := bson.M{"user_id": userId}
	_, err := s.getSessionsCollection().DeleteMany(ctx, filter)
	return err
}
