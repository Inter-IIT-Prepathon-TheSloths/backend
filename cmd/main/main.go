package main

import (
	"context"
	"log"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/database"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/server"
)

func main() {
	client := database.New()
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	srv := server.NewServer(client)

	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
