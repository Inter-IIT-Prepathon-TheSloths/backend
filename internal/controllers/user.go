package controllers

import (
	"fmt"
	"net/http"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/services"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{service: services.NewUserService(client)}
}

func (uc *UserController) SignUp(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return err
	}

	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Validation failed for field '%s', violating '%s' condition\n", err.Field(), err.Tag()))
		}
	}

	existingUser, err := uc.service.GetUser(c.Request().Context(), bson.M{"email": user.Email})
	if err != nil {
		return err
	}

	if existingUser != nil {
		return echo.NewHTTPError(http.StatusConflict, "Email already registered")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	id, err := uc.service.CreateUser(c.Request().Context(), &user)
	if err != nil {
		return err
	}

	jwt, err := utils.CreateJwtToken(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to sign the jwt")
	}

	return c.JSON(http.StatusCreated, map[string]string{"jwt": jwt})
}

func (uc *UserController) GetUsers(c echo.Context) error {
	users, err := uc.service.GetAllUsers(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}
