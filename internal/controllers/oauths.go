package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (uc *UserController) AuthController(c echo.Context) error {
	provider := c.Param("provider")
	state, err := utils.GenerateState()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate state")
	}
	var url string
	if provider == "google" {
		url = config.Google_conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	} else if provider == "github" {
		url = config.Github_conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid provider")
	}
	return c.JSON(http.StatusOK, map[string]string{"url": url})
}

func (uc *UserController) Callback(c echo.Context) error {
	provider := c.Param("provider")
	var conf *oauth2.Config
	var api_url string
	if provider == "google" {
		conf = config.Google_conf
		api_url = "https://www.googleapis.com/oauth2/v2/userinfo"
	} else if provider == "github" {
		conf = config.Github_conf
		api_url = "https://api.github.com/user/emails"
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid provider")
	}

	code := c.Request().URL.Query().Get("code")
	t, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to login with Oauth provider")
	}
	client := conf.Client(context.Background(), t)
	resp, err := client.Get(api_url)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to fetch user details")
	}
	defer resp.Body.Close()

	var userDetails models.GoogleUserDetails
	var emails []models.Email

	if provider == "google" {
		err = json.NewDecoder(resp.Body).Decode(&userDetails)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode the user information")
		}
		emails = []models.Email{
			{Email: userDetails.Email, IsVerified: true},
		}
	} else if provider == "github" {
		var githubEmails []struct {
			Email    string `json:"email"`
			Verified bool   `json:"verified"`
			Primary  bool   `json:"primary"`
		}
		err = json.NewDecoder(resp.Body).Decode(&githubEmails)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode the user information")
		}

		suffix := "@users.noreply.github.com"
		for _, email := range githubEmails {
			if email.Verified && !strings.HasSuffix(email.Email, suffix) {
				emails = append(emails, models.Email{
					Email:      email.Email,
					IsVerified: true,
				})
			}
		}
		if len(emails) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Please use a github account with a public verified email")
		}
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode the userinfo body")
	}

	filter := utils.ConstructEmailFilter(emails)
	existingUser, err := uc.service.GetUser(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	baseUrl := os.Getenv("FRONTEND_URL")
	if existingUser == nil {
		user := &models.User{
			Emails:  emails,
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

	newEmails, needsUpdate := utils.UpdateOauthEmails(existingUser.Emails, emails)
	if needsUpdate {
		existingUser.Emails = newEmails
		err = uc.service.UpdateUser(c.Request().Context(), existingUser.ID.Hex(), existingUser)
		if err != nil {
			return err
		}
	}

	jwt, err := utils.CreateJwtToken(existingUser.ID.Hex())
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/backend_redirect?token=%s", baseUrl, jwt))
}
