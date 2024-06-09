package middlewares

import (
    "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/utils"
	"lineblocs.com/api/user"
)

func BasicAuthMiddleware(userStore user.UserStoreInterface) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if userStore.ValidateAccess(username, password) {
			utils.Log(logrus.InfoLevel, "Authentification is successfully passed")
			utils.SetMicroservice(username)
			return true, nil
		}
		return false, nil
    });
}