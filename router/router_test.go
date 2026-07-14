package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("Should return the test response", func(t *testing.T) {

		e := New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		handler := func(c echo.Context) error {
			return c.String(http.StatusOK, "Test Response")
		}

		if assert.NoError(t, handler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "Test Response")
		}
	})

	t.Run("Should recover panic and return 500", func(t *testing.T) {
		e := New()
		e.GET("/panic", func(c echo.Context) error {
			panic("boom")
		})

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
