package controllers

import (
	"net/http"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
)

func (uc *UserController) SendVerificationMail(c echo.Context) error {
	use := c.Param("use")
	email := c.QueryParam("email")

	emailFilter := utils.ConstructEmailFilter([]models.Email{{Email: email}})
	user, err := uc.service.GetUser(c.Request().Context(), emailFilter)
	if err != nil {
		return err
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	if use == "add_email" || use == "update_password" {
		// presentUser := c.Get("user").(*models.User)

	}
	return nil
}
