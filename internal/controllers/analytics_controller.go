package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/models"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/services"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnalyticsController struct {
	service *services.AnalyticsService
}

func NewAnalyticsController(client *mongo.Client) *AnalyticsController {
	return &AnalyticsController{service: services.NewAnalyticsService(client)}
}

func (ac *AnalyticsController) GetCompanies(c echo.Context) error {
	apiUrl := fmt.Sprintf("%s/companies", config.AnalyticsUrl)
	resp, err := ac.service.CallApi(apiUrl)
	return utils.Proxy(c, resp, err)
}

func (ac *AnalyticsController) GetAnalytics(c echo.Context) error {
	user_id := c.Get("_id").(primitive.ObjectID)
	index := c.Param("index")
	analytics, err := ac.service.Get(c.Request().Context(), bson.M{"user_id": user_id, "company_id": index})
	if err != nil {
		return err
	}

	if analytics == nil {
		apiUrl := fmt.Sprintf("%s/analytics?index=%s", config.AnalyticsUrl, index)
		resp, err := ac.service.CallApi(apiUrl)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var result map[string]interface{}
		if err = json.Unmarshal(body, &result); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse external API response: %w", err)
		}

		compressedBody, err := utils.Compress(body)
		if err != nil {
			return err
		}

		analyticsDoc := models.Analytics{
			UserId:    user_id,
			CompanyId: index,
			Data:      compressedBody,
		}

		err = ac.service.Create(c.Request().Context(), analyticsDoc)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, result)
	}

	decompressedBody, err := utils.Decompress(analytics[0].Data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, decompressedBody)
}

func (ac *AnalyticsController) GetMyAnalytics(c echo.Context) error {
	user_id := c.Get("_id").(primitive.ObjectID)
	analytics, err := ac.service.Get(c.Request().Context(), bson.M{"user_id": user_id})
	if err != nil {
		return err
	}

	fmt.Println(analytics)

	var resultAnalytics []interface{}
	for _, a := range analytics {
		decompressedBody, err := utils.Decompress(a.Data)
		if err != nil {
			return err
		}

		resultAnalytics = append(resultAnalytics, decompressedBody)
	}

	return c.JSON(http.StatusOK, resultAnalytics)
}
