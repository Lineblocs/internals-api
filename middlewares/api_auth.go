package middlewares

import (
	"net/http"
    "github.com/labstack/echo/v4"
)

func APIAuthMiddleware(expectedValue string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            headerKey := "x-lineblocs-api-token"
            header := c.Request().Header.Get(headerKey)
            if header == "" {
                return c.String(http.StatusBadRequest, "Missing header: " + headerKey)
            }
            if header != expectedValue {
                return c.String(http.StatusUnauthorized, "API key is invalid.")
            }
            return next(c)
        }
    }
}