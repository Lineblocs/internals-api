package middlewares

import (
	"net/http"
    "github.com/labstack/echo/v4"
    "github.com/sirupsen/logrus"
	"lineblocs.com/api/utils"
)

func APIAuthMiddleware(expectedValue string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            headerKey := "x-lineblocs-api-token"
            header := c.Request().Header.Get(headerKey)
            if header == "" {
	            utils.Log(logrus.InfoLevel, "API token header is missing")
                return c.String(http.StatusBadRequest, "Missing header: " + headerKey)
            }

	        utils.Log(logrus.InfoLevel, "received API token: '" + header + "'")
	        utils.Log(logrus.InfoLevel, "comparing agaisnt configured token: '" + expectedValue + "'")
            if header != expectedValue {
                return c.String(http.StatusUnauthorized, "API key is invalid.")
            }
            return next(c)
        }
    }
}