package routes

import (
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/controllers"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/middlewares"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUserRoutes(e *echo.Group, client *mongo.Client) {
	userController := controllers.NewUserController(client)

	e.POST("/signup", userController.SignUp)
	e.Use(middlewares.AuthenticationMiddleware)
	e.GET("/", userController.GetUsers)
}
