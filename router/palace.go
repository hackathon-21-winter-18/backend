package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-21-winter-18/backend/model"
	"github.com/labstack/echo/v4"
)

type PalaceRequest struct {
	Name        string       `json:"name"`
	Image       string       `json:"image"`
	EmbededPins []EmbededPin `json:"embededPins"`
}

type EmbededPin struct {
	Number int     `json:"number"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Word   string  `json:"word"`
	Memo   string  `json:"memo"`
}

func postPalace(c echo.Context) error {
	var req PalaceRequest
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	palaceID, err := model.CreatePalace(ctx, userID, req.Name, req.Image)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	for _, embededPin := range req.EmbededPins {
		err = model.CreateEmbededPin(ctx, embededPin.Number, *palaceID, embededPin.X, embededPin.Y, embededPin.Word, embededPin.Memo)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	return nil
}
