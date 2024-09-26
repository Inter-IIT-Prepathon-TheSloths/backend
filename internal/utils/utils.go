package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func CheckHealth(client *mongo.Client, t int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
	defer cancel()

	err := client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
		return err
	}

	return nil
}
