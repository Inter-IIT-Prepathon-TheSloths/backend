package routes

import (
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/controllers"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/middlewares"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterAnalyticsRoutes(e *echo.Group, client *mongo.Client, userController *controllers.UserController) {
	analyticsController := controllers.NewAnalyticsController(client)

	e.Use(middlewares.AuthenticationMiddleware(userController, false))
	e.Use(middlewares.TwofaMiddleware(userController, false))
	e.GET("/companies", analyticsController.GetCompanies)
	e.GET("/:index", analyticsController.GetAnalytics)
	// e.POST("/login", analyticsController.Login)
}
