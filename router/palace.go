package router

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/hackathon-21-winter-18/backend/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type PostPalace struct {
	Name        *string            `json:"name,omitempty"`
	Image       string             `json:"image"`
	EmbededPins []model.EmbededPin `json:"embededPins"`
	CreatedBy   *uuid.UUID         `json:"createdBy,omitempty"`
}
type PutPalace struct {
	Name        *string            `json:"name"`
	Image       string             `json:"image"`
	EmbededPins []model.EmbededPin `json:"embededPins"`
}

type Share struct {
	Share bool `json:"share"`
}

type ID struct {
	ID uuid.UUID `json:"id"`
}

func getPalaces(c echo.Context) error {
	ctx := c.Request().Context()
	palaces, err := model.GetSharePalaces(ctx)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, palace := range palaces {
		palacePins, err := model.GetEmbededPins(ctx, palace.ID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		for _, palacePin := range palacePins {
			palace.EmbededPins = append(palace.EmbededPins, palacePin)
		}

		palace.Image, err = model.EncodeToBase64(ctx, palace.Image)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	// first pin sort
	min := c.QueryParam("minpins")
	max := c.QueryParam("maxpins")
	if min != "" && max != "" && min > max {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	palaces = model.ExtractFromPalacesBasedOnEmbededPins(palaces, max, min)
	// second sort with query
	sortmethod := c.QueryParam("sort")
	switch sortmethod {
	case "first_shared_at":
		fmt.Println("first_shared_at")
		sort.Slice(palaces, func(i, j int) bool {
			return palaces[i].FirstSharedAt.Before(palaces[j].FirstSharedAt)
		})
	case "shared_at":
		sort.Slice(palaces, func(i, j int) bool {
			return palaces[i].SharedAt.Before(palaces[j].SharedAt)
		})
	}

	return echo.NewHTTPError(http.StatusOK, palaces)
}

func getMyPalaces(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		c.Logger().Error(err)
		return errSessionNotFound(err)
	}
	userID, err := uuid.Parse(sess.Values["userID"].(string))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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

		palace.Image, err = model.EncodeToBase64(ctx, palace.Image)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return echo.NewHTTPError(http.StatusOK, palaces)
}

func postPalace(c echo.Context) error {
	var req PostPalace
	sess, err := session.Get("sessions", c)
	if err != nil {
		c.Logger().Error(err)
		return errSessionNotFound(err)
	}
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return errBind(err)
	}
	userID, err := uuid.Parse(sess.Values["userID"].(string))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	ctx := c.Request().Context()
	path, err := model.CreatePathName(ctx, req.Image)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	for _, embededPin := range req.EmbededPins {
		if embededPin.Number == nil || embededPin.X == nil || embededPin.Y == nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid pins"))
		}
	}

	palaceID, err := model.CreatePalace(ctx, userID, req.CreatedBy, req.Name, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = model.DecodeToImageAndSave(ctx, req.Image, path)
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

	res := ID{ID: *palaceID}

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
	sess, err := session.Get("sessions", c)
	if err != nil {
		c.Logger().Error(err)
		return errSessionNotFound(err)
	}
	userID, err := uuid.Parse(sess.Values["userID"].(string))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	ctx := c.Request().Context()
	err = model.CheckPalaceHeldBy(ctx, userID, palaceID)
	if err != nil {
		c.Logger().Error(err)
		generateEchoError(err)
	}
	path, err := model.CreatePathName(ctx, req.Image)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	for _, embededPin := range req.EmbededPins {
		if embededPin.Number == nil || embededPin.X == nil || embededPin.Y == nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid pins"))
		}
	}

	unupdatedPath, err := model.GetPalaceImagePath(ctx, palaceID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.UpdatePalace(ctx, palaceID, req.Name, path)
	if err != nil {
		c.Logger().Error(err)
		return generateEchoError(err)
	}
	err = model.RemoveImage(ctx, unupdatedPath)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.DecodeToImageAndSave(ctx, req.Image, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = model.DeleteEmbededPins(ctx, palaceID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	for _, updatedEmbededPin := range req.EmbededPins {
		err = model.CreateEmbededPin(ctx, updatedEmbededPin.Number, palaceID, updatedEmbededPin.X, updatedEmbededPin.Y, updatedEmbededPin.Word, updatedEmbededPin.Place, updatedEmbededPin.Do)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return echo.NewHTTPError(http.StatusOK)
}

func deletePalace(c echo.Context) error {
	palaceID, err := uuid.Parse(c.Param("palaceID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	sess, err := session.Get("sessions", c)
	if err != nil {
		c.Logger().Error(err)
		return errSessionNotFound(err)
	}
	userID, err := uuid.Parse(sess.Values["userID"].(string))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	ctx := c.Request().Context()
	err = model.CheckPalaceHeldBy(ctx, userID, palaceID)
	if err != nil {
		c.Logger().Error(err)
		generateEchoError(err)
	}

	unupdatedPath, err := model.GetPalaceImagePath(ctx, palaceID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.DeletePalace(ctx, palaceID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.RemoveImage(ctx, unupdatedPath)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK)
}
