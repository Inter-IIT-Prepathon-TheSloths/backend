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

type UserService struct {
	client *mongo.Client
}

func NewUserService(client *mongo.Client) *UserService {
	return &UserService{
		client: client,
	}
}

func (s *UserService) getUserCollection() *mongo.Collection {
	return s.client.Database(config.DbName).Collection("users")
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (string, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	insertedDoc, err := s.getUserCollection().InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	objectId, ok := insertedDoc.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("InsertedID is not of type primitive.ObjectID")
	}
	return objectId.Hex(), err
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	cur, err := s.getUserCollection().Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []models.User

	for cur.Next(ctx) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetUser(ctx context.Context, filter bson.M) (*models.User, error) {
	var user models.User
	err := s.getUserCollection().FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error Decoding User Doc")
	}

	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user *models.User) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": user}
	_, err = s.getUserCollection().UpdateOne(ctx, filter, update)

	return err
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	_, err = s.getUserCollection().DeleteOne(ctx, filter)
	return err
}
