package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

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

	baseUrl := os.Getenv("FRONTEND_URL")
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
