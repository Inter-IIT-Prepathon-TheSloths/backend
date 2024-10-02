package services

import (
	"context"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *UserService) getTwofaCollection() *mongo.Collection {
	return s.client.Database(config.DbName).Collection("twofa")
}

func (s *UserService) UpdateTwofa(ctx context.Context, user_id primitive.ObjectID, twofa *models.TwoFactor) error {
	filter := bson.M{"user_id": user_id}
	update := bson.M{"$set": twofa}

	opts := options.Update().SetUpsert(true)
	_, err := s.getTwofaCollection().UpdateOne(ctx, filter, update, opts)
	return err
}

func (s *UserService) DeleteTwoFactor(ctx context.Context, user_id primitive.ObjectID) error {
	filter := bson.M{"user_id": user_id}
	_, err := s.getTwofaCollection().DeleteOne(ctx, filter)
	return err
}

func (s *UserService) GetTwoFactor(ctx context.Context, user_id primitive.ObjectID) (*models.TwoFactor, error) {
	var twofa models.TwoFactor
	err := s.getTwofaCollection().FindOne(ctx, bson.M{"user_id": user_id}).Decode(&twofa)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &twofa, nil
}
