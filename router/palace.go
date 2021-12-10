package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-21-winter-18/backend/model"
	"github.com/labstack/echo/v4"
)

type PostPalace struct {
	Name        *string             `json:"name,omitempty"`
	Image       string             `json:"image"`
	EmbededPins []model.EmbededPin `json:"embededPins"`
	CreatedBy   *uuid.UUID          `json:"createdBy,omitempty"`
}
type PutPalace struct {
	Name        string             `json:"name"`
	Image       string             `json:"image"`
	EmbededPins []model.EmbededPin `json:"embededPins"`
}

type Share struct {
	Share bool `json:"share"`
}

type P struct {
	//TODO 多分消す
	ID uuid.UUID `json:"id"`
}

func getPalaces(c echo.Context) error {
	return nil
}

func getMyPalaces(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	palaces, err := model.GetPalaces(ctx, userID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, palace := range palaces {
		embededPins, err := model.GetEmbededPins(ctx, palace.ID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		for _, embededPin := range embededPins {
			palace.EmbededPins = append(palace.EmbededPins, embededPin)
		}

		palace.Image, err = model.EncodeTobase64(ctx, palace.Image)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return echo.NewHTTPError(http.StatusOK, palaces)
}

func postPalace(c echo.Context) error {
	var req PostPalace
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return errBind(err)
	}

	ctx := c.Request().Context()
	path, err := model.CreatePathName(ctx, req.Image)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	palaceID, err := model.CreatePalace(ctx, userID, req.CreatedBy, req.Name, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	for _, embededPin := range req.EmbededPins {
		err = model.CreateEmbededPin(ctx, embededPin.Number, *palaceID, embededPin.X, embededPin.Y, embededPin.Word, embededPin.Place, embededPin.Do)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	//TODO model関数この順番でいいのか
	err = model.DecodeToImageAndSave(ctx, req.Image, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	//TODO 多分resけす
	res := P{ID: *palaceID}

	return echo.NewHTTPError(http.StatusOK, res)
}

func putPalace(c echo.Context) error {
	var req PutPalace
	palaceID, err := uuid.Parse(c.Param("palaceID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	unupdatedPath, err := model.GetPalaceImagePath(ctx, palaceID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = model.RemoveImage(ctx, unupdatedPath)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	path, err := model.CreatePathName(ctx, req.Image)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = model.UpdatePalace(ctx, palaceID, req.Name, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = model.DeleteEmbededPins(ctx, palaceID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	
	for _, updatedEmbededPin := range req.EmbededPins {
		err = model.CreateEmbededPin(ctx, updatedEmbededPin.Number, palaceID, updatedEmbededPin.X, updatedEmbededPin.Y, updatedEmbededPin.Word, updatedEmbededPin.Place, updatedEmbededPin.Do)
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

func deletePalace(c echo.Context) error {
	palaceID, err := uuid.Parse(c.Param("palaceID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	unupdatedPath, err := model.GetPalaceImagePath(ctx, palaceID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.RemoveImage(ctx, unupdatedPath)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = model.DeletePalace(ctx, palaceID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return echo.NewHTTPError(http.StatusOK)
}

func sharePalace(c echo.Context) error {
	var req Share
	palaceID, err := uuid.Parse(c.Param("palaceID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	err = model.SharePalace(ctx, palaceID, req.Share)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return echo.NewHTTPError(http.StatusOK)
}
