package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func errSessionNotFound(err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Failed in Getting Session:%w", err).Error())
}
