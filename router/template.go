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

type PostTemplate struct {
	Name         *string             `json:"name,omitempty"`
	Image        string              `json:"image"`
	TemplatePins []model.TemplatePin `json:"pins"`
	CreatedBy    *uuid.UUID          `json:"createdBy,omitempty"`
}

type PutTemplate struct {
	Name         *string             `json:"name"`
	Image        string              `json:"image"`
	TemplatePins []model.TemplatePin `json:"pins"`
}

func getTemplates(c echo.Context) error {
	ctx := c.Request().Context()
	templates, err := model.GetShareTemplates(ctx)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	for _, template := range templates {
		templatePins, err := model.GetTemplatePins(ctx, template.ID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		for _, templatePin := range templatePins {
			template.TemplatePins = append(template.TemplatePins, templatePin)
		}
		template.Image, err = model.EncodeToBase64(ctx, template.Image)
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
	templates = model.ExtractFromTemplatesBasedOnTemplatePins(templates, max, min)
	// second sort with query
	sortmethod := c.QueryParam("sort")
	switch sortmethod {
	case "first_shared_at":
		sort.Slice(templates, func(i, j int) bool {
			return templates[i].FirstSharedAt.Before(templates[j].FirstSharedAt)
		})
	case "shared_at":
		sort.Slice(templates, func(i, j int) bool {
			return templates[i].SharedAt.Before(templates[j].SharedAt)
		})
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

	ctx := c.Request().Context()
	templates, err := model.GetTemplates(ctx, userID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, template := range templates {
		templatePins, err := model.GetTemplatePins(ctx, template.ID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		for _, templatePin := range templatePins {
			template.TemplatePins = append(template.TemplatePins, templatePin)
		}

		template.Image, err = model.EncodeToBase64(ctx, template.Image)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return echo.NewHTTPError(http.StatusOK, templates)
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
	for _, templatePin := range req.TemplatePins {
		if templatePin.Number == nil || templatePin.X == nil || templatePin.Y == nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid pins"))
		}
	}

	templateID, err := model.CreateTemplate(ctx, userID, req.CreatedBy, req.Name, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.DecodeToImageAndSave(ctx, req.Image, path)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	for _, templatePin := range req.TemplatePins {
		err = model.CreateTemplatePin(ctx, templatePin.Number, *templateID, templatePin.X, templatePin.Y)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
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
	for _, templatePin := range req.TemplatePins {
		if templatePin.Number == nil || templatePin.X == nil || templatePin.Y == nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid pins"))
		}
	}

	unupdatedPath, err := model.GetTemplateImagePath(ctx, templateID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	err = model.UpdateTemplate(ctx, templateID, req.Name, path)
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

	err = model.DeleteTemplatePins(ctx, templateID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	for _, updatedTemplatePin := range req.TemplatePins {
		err = model.CreateTemplatePin(ctx, updatedTemplatePin.Number, templateID, updatedTemplatePin.X, updatedTemplatePin.Y)
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
