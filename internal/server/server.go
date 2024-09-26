package server

import (
	"log"
	"net/http"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}

	log.Printf("Error: %v", err)
	c.JSON(code, map[string]string{"error": msg})
}

func NewServer(client *mongo.Client) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register user routes
	e.GET("/", HealthCheck(client))

	e.HTTPErrorHandler = customHTTPErrorHandler

	return e
}

func HealthCheck(client *mongo.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		utils.CheckHealth(client, 1)
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	}
}
