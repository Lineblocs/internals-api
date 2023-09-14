package handler

import (
	"os"
	"testing"

	helpers "github.com/Lineblocs/go-helpers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	helpers.InitLogrus("stdout")

	t.Run("Should Register routes", func(t *testing.T) {
		e := echo.New()
		os.Setenv("USE_AUTH_MIDDLEWARE", "on")

		h := &Handler{}
		h.Register(e)

		assert.NotNil(t, h)
		assert.NotNil(t, e)
		assert.NotEmpty(t, e)
	})
}
