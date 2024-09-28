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
	e.POST("/login", userController.Login)
	e.GET("/google", userController.GoogleAuthController)
	e.GET("/callback/google", userController.CallbackGoogle)
	e.POST("/create_password", userController.CreatePassword)
	e.GET("/verify_email", userController.VerifyEmail)

	e.Use(middlewares.AuthenticationMiddleware)
	e.GET("/me", userController.GetMyDetails)
	e.POST("/add_email", userController.AddEmail)
	e.GET("/send_activation", userController.SendActivationMail)
}
