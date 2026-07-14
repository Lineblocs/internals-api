package router

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func RecoveryMiddleware() echo.MiddlewareFunc {
	config := middleware.RecoverConfig{
		StackSize:         8 << 10,
		DisableStackAll:   true,
		DisablePrintStack: true,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logrus.WithFields(logrus.Fields{
				"path":   c.Request().URL.Path,
				"method": c.Request().Method,
				"ip":     c.RealIP(),
				"stack":  string(stack),
			}).WithError(err).Error("panic recovered in HTTP handler")
			return nil
		},
	}

	return middleware.RecoverWithConfig(config)
}

func PanicRouteHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("%s", http.StatusText(http.StatusInternalServerError)),
		})
	}
}
