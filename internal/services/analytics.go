package services

import (
	"context"
	"net/http"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnalyticsService struct {
	client *mongo.Client
}

func NewAnalyticsService(client *mongo.Client) *AnalyticsService {
	return &AnalyticsService{
		client: client,
	}
}

func (s *AnalyticsService) getCollection() *mongo.Collection {
	return s.client.Database(config.DbName).Collection("analytics")
}

func (as *AnalyticsService) CallApi(apiUrl string) (*http.Response, error) {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to call external API")
	}

	return resp, nil
}

func (as *AnalyticsService) Create(ctx context.Context, analytics models.Analytics) error {
	_, err := as.getCollection().InsertOne(ctx, analytics)
	return err
}

func (as *AnalyticsService) Get(ctx context.Context, user_id primitive.ObjectID, index string) (*models.Analytics, error) {
	var analytics models.Analytics
	err := as.getCollection().FindOne(ctx, bson.M{"user_id": user_id, "company_id": index}).Decode(&analytics)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &analytics, nil
}

func (as *AnalyticsService) Delete(ctx context.Context, user_id primitive.ObjectID, company_id string) error {
	_, err := as.getCollection().DeleteOne(ctx, bson.M{"user_id": user_id, "company_id": company_id})
	if err != nil {
		return err
	}

	return nil
}
