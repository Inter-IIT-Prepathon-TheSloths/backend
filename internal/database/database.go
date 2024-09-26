package database

import (
	"context"
	"log"
	"os"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db_url = os.Getenv("DB_URL")
)

func New() *mongo.Client {
	clientOptions := options.Client().ApplyURI(db_url)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	utils.CheckHealth(client, 5)

	return client
}
