package database

import (
	"context"
	"fmt"
	"log"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db_url = config.DbUrl
)

func createIndex(client *mongo.Client) {
	collection := client.Database(config.DbName).Collection("verifications")

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"expires_at": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	indexName, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Index created: ", indexName)
}

func New() *mongo.Client {
	clientOptions := options.Client().ApplyURI(db_url)
	client, err := mongo.Connect(context.Background(), clientOptions)

	createIndex(client)

	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	if err := utils.CheckHealth(client, 5); err != nil {
		log.Fatalf("db not accessible: %v", err)
	}

	return client
}
