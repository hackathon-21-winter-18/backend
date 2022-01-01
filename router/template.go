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

type PostTemplate struct {
	Name       *string     `json:"name,omitempty"`
	Image      string      `json:"image"`
	Pins       []model.Pin `json:"pins"`
	CreatedBy  *uuid.UUID  `json:"createdBy,omitempty"`
	OriginalID *uuid.UUID  `json:"originalID"`
}

type PutTemplate struct {
	Name  *string     `json:"name"`
	Image string      `json:"image"`
	Pins  []model.Pin `json:"pins"`
}

func getSharedTemplates(c echo.Context) error {
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
		Sort:           sort,
		MaxEmbededPins: max,
		MinEmbededPins: min,
	}

	ctx := c.Request().Context()
	templates, err := model.GetSharedTemplates(ctx, requestQuery)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	for _, template := range templates {
		pins, err := model.GetPins(ctx, template.ID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		for _, Pin := range pins {
			template.Pins = append(template.Pins, Pin)
		}
		template.Image, err = model.EncodeToBase64(ctx, template.Image)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return echo.NewHTTPError(http.StatusOK, templates)
}

func getMyTemplates(c echo.Context) error {
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
	if sort != "" && sort != "updated_at" && sort != "-updated_at" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid sort query"))
	}

	requestQuery := model.RequestQuery{
		Sort: sort,
		MaxEmbededPins: max,
		MinEmbededPins: min,
	}

	ctx := c.Request().Context()
	templates, err := model.GetMyTemplates(ctx, userID, requestQuery)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, template := range templates {
		pins, err := model.GetPins(ctx, template.ID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		for _, pin := range pins {
			template.Pins = append(template.Pins, pin)
		}

		template.Image, err = model.EncodeToBase64(ctx, template.Image)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return echo.NewHTTPError(http.StatusOK, templates)
}

func getTemplate(c echo.Context) error {
	templateID, err := uuid.Parse(c.Param("templateID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	template, err := model.GetTemplate(ctx, templateID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	pins, err := model.GetPins(ctx, template.ID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	for _, pin := range pins {
		template.Pins = append(template.Pins, pin)
	}

	template.Image, err = model.EncodeToBase64(ctx, template.Image)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK, template)
}

func postTemplate(c echo.Context) error {
	var req PostTemplate
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
	for _, pin := range req.Pins {
		if pin.Number == nil || pin.X == nil || pin.Y == nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid pins"))
		}
	}
	number_of_pins := len(req.Pins)

	templateID, err := model.CreateTemplate(ctx, req.OriginalID, userID, req.CreatedBy, req.Name, number_of_pins, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.DecodeToImageAndSave(ctx, req.Image, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, pin := range req.Pins {
		err = model.CreatePin(ctx, pin.Number, *templateID, pin.X, pin.Y)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if *req.CreatedBy != userID {
		err = model.RecordTemplateSavingUser(ctx, *req.OriginalID, userID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	res := ID{ID: *templateID}

	return echo.NewHTTPError(http.StatusOK, res)
}

func putTemplate(c echo.Context) error {
	var req PutTemplate
	templateID, err := uuid.Parse(c.Param("templateID"))
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
	err = model.CheckTemplateHeldBy(ctx, userID, templateID)
	if err != nil {
		c.Logger().Error(err)
		return generateEchoError(err)
	}
	path, err := model.CreatePathName(ctx, req.Image)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	for _, pin := range req.Pins {
		if pin.Number == nil || pin.X == nil || pin.Y == nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid pins"))
		}
	}
	number_of_pins := len(req.Pins)

	unupdatedPath, err := model.GetTemplateImagePath(ctx, templateID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.UpdateTemplate(ctx, templateID, req.Name, number_of_pins, path)
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

	err = model.DeletePins(ctx, templateID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	for _, updatedPin := range req.Pins {
		err = model.CreatePin(ctx, updatedPin.Number, templateID, updatedPin.X, updatedPin.Y)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return echo.NewHTTPError(http.StatusOK)
}

func deleteTemplate(c echo.Context) error {
	templateID, err := uuid.Parse(c.Param("templateID"))
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
	err = model.CheckTemplateHeldBy(ctx, userID, templateID)
	if err != nil {
		c.Logger().Error(err)
		return generateEchoError(err)
	}

	unupdatedPath, err := model.GetTemplateImagePath(ctx, templateID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.DeleteTemplate(ctx, templateID)
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

func shareTemplate(c echo.Context) error {
	var req Share
	templateID, err := uuid.Parse(c.Param("templateID"))
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
	err = model.CheckTemplateHeldBy(ctx, userID, templateID)
	if err != nil {
		c.Logger().Error(err)
		return generateEchoError(err)
	}

	err = model.ShareTemplate(ctx, templateID, req.Share)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return echo.NewHTTPError(http.StatusOK)
}
