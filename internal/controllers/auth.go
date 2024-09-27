package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/services"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
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

	validation_err := utils.Validate(user)
	if validation_err != nil {
		return validation_err
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

func (uc *UserController) Login(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return err
	}

	existingUser, err := uc.service.GetUser(c.Request().Context(), bson.M{"email": user.Email})
	if err != nil {
		return err
	}

	if existingUser == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	err = utils.VerifyPassword(existingUser.Password, user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	jwt, err := utils.CreateJwtToken(existingUser.ID.Hex())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"jwt": jwt})
}

func (uc *UserController) GoogleAuthController(c echo.Context) error {
	state, err := utils.GenerateState()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate state")
	}
	url := config.Google_conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return c.JSON(http.StatusOK, url)
}

func (uc *UserController) CallbackGoogle(c echo.Context) error {
	conf := config.Google_conf
	code := c.Request().URL.Query().Get("code")
	fmt.Println(code)
	t, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to login with google")
	}
	client := conf.Client(context.Background(), t)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to fetch user details")
	}
	defer resp.Body.Close()

	var userDetails models.UserDetails

	err = json.NewDecoder(resp.Body).Decode(&userDetails)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode the userinfo body")
	}

	existingUser, err := uc.service.GetUser(c.Request().Context(), bson.M{"email": userDetails.Email})
	if err != nil {
		return err
	}

	if existingUser == nil {
		user := &models.User{
			Email:    userDetails.Email,
			Picture:  userDetails.Picture,
			Username: userDetails.Name,
		}
		idCreated, err := uc.service.CreateUser(c.Request().Context(), user)
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?id=%s", config.Frontend_password, idCreated))
	}

	if existingUser.Password == "" {
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?id=%s", config.Frontend_password, existingUser.ID))
	}

	jwt, err := utils.CreateJwtToken(existingUser.ID.Hex())
	if err != nil {
		return err
	}

	fmt.Printf("User: %v", userDetails)
	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?token=%s", config.Frontend_home, jwt))
}

func (uc *UserController) CreatePassword(c echo.Context) error {
	var password models.UserPassword
	if err := c.Bind(&password); err != nil {
		return err
	}

	validation_error := utils.Validate(password)
	if validation_error != nil {
		return validation_error
	}

	user, err := uc.service.GetUser(c.Request().Context(), bson.M{"_id": password.ID})
	if err != nil {
		return err
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	hashedPassword, err := utils.HashPassword(password.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	err = uc.service.UpdateUser(c.Request().Context(), password.ID, user)
	if err != nil {
		return err
	}

	jwt, err := utils.CreateJwtToken(user.ID.Hex())
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?token=%s", config.Frontend_home, jwt))
}

func (uc *UserController) GetUsers(c echo.Context) error {
	users, err := uc.service.GetAllUsers(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}
