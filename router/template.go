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

	//TODO 多分resけす
	res := P{ID: *templateID}

	return echo.NewHTTPError(http.StatusOK, res)
}
