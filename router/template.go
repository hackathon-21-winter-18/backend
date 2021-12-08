package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-21-winter-18/backend/model"
	"github.com/labstack/echo/v4"
)

type PostTemplate struct {
	Name         string              `json:"name"`
	Image        string              `json:"image"`
	TemplatePins []model.TemplatePin `json:"templatePins"`
	CreatedBy    uuid.UUID           `json:"createdBy"`
}

func getTemplates(c echo.Context) error {
	ctx := c.Request().Context()
	templates, err := model.GetTemplates(ctx)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	for _, template := range templates {
		templatePins, err := model.GetTemplatePins(ctx, template.ID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		for _, templatePin := range templatePins {
			template.TemplatePins = append(template.TemplatePins, templatePin)
		}

		template.Image, err = model.EncodeTobase64(ctx, template.Image)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	return echo.NewHTTPError(http.StatusOK, templates)
}

func postTemplate(c echo.Context) error {
	var req PostTemplate
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	path, err := model.CreatePathName(ctx, req.Image)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	templateID, err := model.CreateTemplate(ctx, userID, req.CreatedBy, req.Name, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	for _, templatePin := range req.TemplatePins {
		// create template pin
		err = model.CreateTemplatePin(ctx, templatePin.Number, *templateID, templatePin.X, templatePin.Y)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	err = model.DecodeToImageAndSave(ctx, req.Image, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return echo.NewHTTPError(http.StatusOK)
}
