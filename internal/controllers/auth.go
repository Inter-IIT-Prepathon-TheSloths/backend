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
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var userDetails models.LoginUserDetails

	if err := c.Bind(&userDetails); err != nil {
		return err
	}

	validation_err := utils.Validate(userDetails)
	if validation_err != nil {
		return validation_err
	}

	filter := utils.ConstructEmailFilter(userDetails.Email)
	existingUser, err := uc.service.GetUser(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Email already registered")
	}

	hashedPassword, err := utils.HashPassword(userDetails.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Emails:   []models.Email{{Email: userDetails.Email, IsVerified: false}},
		Password: hashedPassword,
	}

	id, err := uc.service.CreateUser(c.Request().Context(), user)
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
	var userDetails models.LoginUserDetails
	if err := c.Bind(&userDetails); err != nil {
		return err
	}

	filter := utils.ConstructEmailFilter(userDetails.Email)
	existingUser, err := uc.service.GetUser(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	emailBody := utils.GetEmailBody(userDetails.Email, existingUser.Emails)
	if !emailBody.IsVerified {
		return echo.NewHTTPError(http.StatusBadRequest, "please verify your email to continue")
	}

	if existingUser.Password == "" {
		baseUrl := c.Request().Proto + "://" + c.Request().Host
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?id=%s", fmt.Sprintf("%s/backend_redirect", baseUrl), existingUser.ID.Hex()))
	}

	err = utils.VerifyPassword(existingUser.Password, userDetails.Password)
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
	return c.JSON(http.StatusOK, map[string]string{"url": url})
}

func (uc *UserController) CallbackGoogle(c echo.Context) error {
	conf := config.Google_conf
	code := c.Request().URL.Query().Get("code")
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

	var userDetails models.GoogleUserDetails

	err = json.NewDecoder(resp.Body).Decode(&userDetails)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode the userinfo body")
	}

	filter := utils.ConstructEmailFilter(userDetails.Email)
	existingUser, err := uc.service.GetUser(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	baseUrl := c.Request().Proto + "://" + c.Request().Host
	if existingUser == nil {
		user := &models.User{
			Emails:  []models.Email{{Email: userDetails.Email, IsVerified: true}},
			Picture: userDetails.Picture,
			Name:    userDetails.Name,
		}
		idCreated, err := uc.service.CreateUser(c.Request().Context(), user)
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/backend_redirect?id=%s", baseUrl, idCreated))
	}

	if existingUser.Password == "" {
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/backend_redirect?id=%s", baseUrl, existingUser.ID.Hex()))
	}

	jwt, err := utils.CreateJwtToken(existingUser.ID.Hex())
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/backend_redirect?token=%s", baseUrl, jwt))
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

	oid, err := primitive.ObjectIDFromHex(password.ID)
	if err != nil {
		return err
	}

	user, err := uc.service.GetUser(c.Request().Context(), bson.M{"_id": oid})
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

	return c.JSON(http.StatusCreated, map[string]string{"token": jwt})
}

func (uc *UserController) GetMyDetails(c echo.Context) error {
	id := c.Get("_id").(string)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	user, err := uc.service.GetUser(c.Request().Context(), bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (uc *UserController) AddEmail(c echo.Context) error {
	var email models.Email
	if err := c.Bind(&email); err != nil {
		return err
	}

	validation_error := utils.Validate(email)
	if validation_error != nil {
		return validation_error
	}

	id := c.Get("_id").(string)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	user, err := uc.service.GetUser(c.Request().Context(), bson.M{"_id": oid})
	if err != nil {
		return err
	}

	emailExisting := utils.GetEmailBody(email.Email, user.Emails)
	if emailExisting.Email != "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email already exists")
	}

	if len(user.Emails) >= 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "Maximum number of emails reached")
	}

	email.IsVerified = false
	user.Emails = append(user.Emails, email)

	err = uc.service.UpdateUser(c.Request().Context(), id, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Email added successfully"})
}
