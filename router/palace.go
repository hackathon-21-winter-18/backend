package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/hackathon-21-winter-18/backend/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// type PalaceResponse struct {
// 	ID          uuid.UUID          `json:"id"`
// 	Name        string             `json:"name"`
// 	Image       string             `json:"image"`
// 	EmbededPins []model.EmbededPin `json:"embededPins"`
// 	Share       bool               `json:"share"`
// 	SavedCount  int                `json:"savedCount"`
// }

type PostPalace struct {
	Name        *string            `json:"name,omitempty"`
	Image       string             `json:"image"`
	EmbededPins []model.EmbededPin `json:"embededPins"`
	CreatedBy   *uuid.UUID         `json:"createdBy,omitempty"`
	OriginalID  *uuid.UUID         `json:"originalID"`
}
type PutPalace struct {
	Name        *string            `json:"name"`
	Image       string             `json:"image"`
	EmbededPins []model.EmbededPin `json:"embededPins"`
}

func getSharedPalaces(c echo.Context) error {
	var err error
	var max int
	if c.QueryParam("maxpins") != "" {
		max, err = strconv.Atoi(c.QueryParam("maxpins"))
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if max <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("maxpins can't be 0 or negative number"))
		}
	}
	var min int
	if c.QueryParam("minpins") != "" {
		min, err = strconv.Atoi(c.QueryParam("minpins"))
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if min < 0 {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("minpins can't be negative number"))
		}
	}
	if max < min && max != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid pins query"))
	}
	sort := c.QueryParam("sort")
	if sort != "" && sort != "first_shared_at" && sort != "shared_at" && sort != "savedCount" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid sort query"))
	}

	requestQuery := model.RequestQuery{
		Sort: sort,
		MaxEmbededPins: max,
		MinEmbededPins: min,
	}

	ctx := c.Request().Context()
	palaces, err := model.GetSharedPalaces(ctx, requestQuery)
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
	palaces, err := model.GetMyPalaces(ctx, userID)
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

func getPalace(c echo.Context) error {
	palaceID, err := uuid.Parse(c.Param("palaceID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	palace, err := model.GetPalace(ctx, palaceID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

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

	return echo.NewHTTPError(http.StatusOK, palace)
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
	number_of_embededPins := len(req.EmbededPins)

	palaceID, err := model.CreatePalace(ctx, req.OriginalID, userID, req.CreatedBy, req.Name, number_of_embededPins, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.DecodeToImageAndSave(ctx, req.Image, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, embededPin := range req.EmbededPins {
		err = model.CreateEmbededPin(ctx, embededPin.Number, *palaceID, embededPin.X, embededPin.Y, embededPin.Word, embededPin.Place, embededPin.Situation)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if *req.CreatedBy != userID {
		err = model.RecordPalaceSavingUser(ctx, *req.OriginalID, userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
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
		return generateEchoError(err)
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
	number_of_embededPins := len(req.EmbededPins)

	unupdatedPath, err := model.GetPalaceImagePath(ctx, palaceID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.UpdatePalace(ctx, palaceID, req.Name, number_of_embededPins, path)
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
		err = model.CreateEmbededPin(ctx, updatedEmbededPin.Number, palaceID, updatedEmbededPin.X, updatedEmbededPin.Y, updatedEmbededPin.Word, updatedEmbededPin.Place, updatedEmbededPin.Situation)
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
		return generateEchoError(err)
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
		return generateEchoError(err)
	}

	err = model.SharePalace(ctx, palaceID, req.Share)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK)
}
