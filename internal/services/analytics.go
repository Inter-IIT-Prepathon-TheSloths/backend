package services

import (
	"context"
	"log"
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

func (as *AnalyticsService) Get(ctx context.Context, filter bson.M) ([]models.Analytics, error) {
	// var analytics models.Analytics
	cur, err := as.getCollection().Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var analyticsList []models.Analytics

	for cur.Next(ctx) {
		var analytics models.Analytics
		err := cur.Decode(&analytics)
		if err != nil {
			log.Fatal(err)
		}
		analyticsList = append(analyticsList, analytics)
	}

	return analyticsList, nil
}

func (as *AnalyticsService) Delete(ctx context.Context, user_id primitive.ObjectID, company_id string) error {
	_, err := as.getCollection().DeleteOne(ctx, bson.M{"user_id": user_id, "company_id": company_id})
	if err != nil {
		return err
	}

	return nil
}
