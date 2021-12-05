package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-winter-18/backend/model"
	"github.com/labstack/echo/v4"
)

type TemplateRequest struct {
	Name      string      `json:"name"`
	Image     string      `json:"image"`
	Pins      []model.Pin `json:"pins"`
	CreatedBy uuid.UUID   `json:"created_by"`
}

func PostTemplates(c echo.Context) error {
	var req TemplateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	ctx := c.Request().Context()
	_, err := model.CreateTemplate(ctx, req.Name, req.Image, req.Pins, req.CreatedBy)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return nil
}
