package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	helpers "github.com/Lineblocs/go-helpers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"lineblocs.com/api/mocks"
)

func TestVerifyCaller(t *testing.T) {

	mockInstance := Handler{}
	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return error for invalid workspace_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/verifyCaller?workspace_id=abc&number=12345", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := mockInstance.VerifyCaller(c); err != nil {
			t.Errorf("Error: %v", err)
		}
	})
}

func TestCaptureSIPMessage(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/captureSIPMessage?domain=123&sip_msg=12345", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().CaptureSIPMessage("123", "12345").Return(nil, nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if assert.NoError(t, handler.CaptureSIPMessage(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return error when invalid domain or sip_msg", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/captureSIPMessage?domain=123&sip_msg=12345", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().CaptureSIPMessage("123", "12345").Return(nil, errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if err := handler.CaptureSIPMessage(c); err != nil {
			t.Errorf("Error: %v", err)
		}

		if assert.NoError(t, handler.CaptureSIPMessage(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestLogCallInviteEvent(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/logCallInviteEvent", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().LogCallInviteEvent("").Return(nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if assert.NoError(t, handler.LogCallInviteEvent(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/logCallInviteEvent", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().LogCallInviteEvent("").Return(errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if assert.NoError(t, handler.LogCallInviteEvent(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestLogCallByeEvent(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/path", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().LogCallByeEvent("").Return(nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if assert.NoError(t, handler.LogCallByeEvent(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/path", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().LogCallByeEvent("").Return(errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if assert.NoError(t, handler.LogCallByeEvent(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestProcessDialplan(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/processDialplan", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().ProcessDialplan("").Return(nil, nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if assert.NoError(t, handler.ProcessDialplan(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/processDialplan", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().ProcessDialplan("").Return(nil, errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if assert.NoError(t, handler.ProcessDialplan(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestProcessSIPTrunkCall(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/processDialplan", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().ProcessSIPTrunkCall("").Return(nil, nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if assert.NoError(t, handler.ProcessSIPTrunkCall(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/processDialplan", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().ProcessSIPTrunkCall("").Return(nil, errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if assert.NoError(t, handler.ProcessSIPTrunkCall(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestGetUserByDID(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/getUserByDID", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().GetUserByDID("").Return("", nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if assert.NoError(t, handler.GetUserByDID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return error", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/getUserByDID", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockStore := mocks.UserStoreInterface{}
		mockStore.EXPECT().GetUserByDID("").Return("", errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockStore)

		if assert.NoError(t, handler.GetUserByDID(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}
