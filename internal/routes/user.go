package routes

import (
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/controllers"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/middlewares"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUserRoutes(e *echo.Group, client *mongo.Client) {
	userController := controllers.NewUserController(client)

	e.POST("/login", userController.Login)
	e.GET("/oauth/:provider", userController.AuthController)
	e.GET("/callback/:provider", userController.Callback)
	e.POST("/create_password", userController.CreatePassword)

	e.POST("/send_verification/:use", userController.SendVerificationMail)
	e.POST("/verify_code/:use", userController.VerifyVerificationCode)

	e.Use(middlewares.AuthenticationMiddleware(userController))
	e.GET("/me", userController.GetMyDetails)
	e.POST("/add_email", userController.AddEmail)
	e.GET("/generate_2fasecret", userController.Generate2faSecret)

	e.Use(middlewares.TwofaMiddleware(userController))
	e.GET("/regenerate_2fasecret", userController.Generate2faSecret)
	e.DELETE("/disable_2fa", userController.Disable2fa)
	e.GET("/enable_2fa", userController.Enable2fa)
	e.GET("/get_2fa", userController.GetTwofaInfo)
	e.GET("/regenerate_backups_2fa", userController.RegenerateBackup2fa)
}
