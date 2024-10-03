package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
)

func (uc *UserController) Login(c echo.Context) error {
	var userDetails models.LoginUserDetails
	if err := c.Bind(&userDetails); err != nil {
		return err
	}

	filter := utils.ConstructEmailFilter([]models.Email{
		{Email: userDetails.Email, IsVerified: true},
	})
	existingUser, err := uc.service.GetUser(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid Credentials")
	}

	// If the user logged in with Oauth but didn't create password, so redirect him to first create password
	if existingUser.Password == "" {
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?id=%s", fmt.Sprintf("%s/backend_redirect", config.FrontendUrl), existingUser.ID.Hex()))
	}

	err = utils.VerifyPassword(existingUser.Password, userDetails.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	jwt, err := utils.CreateJwtToken(existingUser.ID.Hex(), false, false)
	if err != nil {
		return err
	}

	askForTwofa := existingUser.TwofaEnabled

	return c.JSON(http.StatusOK, map[string]string{"token": jwt, "ask_for_twofa": strconv.FormatBool(askForTwofa)})
}
