package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}
