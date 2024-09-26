package services

import (
	"context"
	"os"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	client *mongo.Client
}

func NewUserService(client *mongo.Client) *UserService {
	return &UserService{client: client}
}

func (s *UserService) getCollection() *mongo.Collection {
	return s.client.Database(os.Getenv("DB_NAME")).Collection("users")
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := s.getCollection().InsertOne(ctx, user)
	return err
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	cur, err := s.getCollection().Find(ctx, nil)
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

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oid}
	var user models.User
	err = s.getCollection().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user *models.User) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	user.ID = oid
	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	_, err = s.getCollection().UpdateOne(ctx, filter, update)
	return err
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	_, err = s.getCollection().DeleteOne(ctx, filter)
	return err
}
