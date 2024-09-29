package controllers

import (
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/services"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{service: services.NewUserService(client)}
}

func (uc *UserController) GetUserService() *services.UserService {
	return uc.service
}
