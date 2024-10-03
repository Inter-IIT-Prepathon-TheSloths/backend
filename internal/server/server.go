package server

import (
	"log"
	"net/http"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/routes"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := "Internal Server Error"

	he, ok := err.(*echo.HTTPError)
	if ok {
		code = he.Code
		msg = he.Message.(string)
	}

	log.Printf("Error: %v", err)
	if err := c.JSON(code, map[string]string{"error": msg}); err != nil {
		log.Printf("Error: Failed to send JSON response: %v", err)
	}
}

func NewServer(client *mongo.Client) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register user routes
	e.GET("/", HealthCheck(client))

	api := e.Group("/api/v1")
	authRouter := api.Group("/auth")

	routes.RegisterUserRoutes(authRouter, client)

	e.HTTPErrorHandler = customHTTPErrorHandler

	return e
}

func HealthCheck(client *mongo.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := utils.CheckHealth(client, 1); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	}
}
