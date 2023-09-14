package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	helpers "github.com/Lineblocs/go-helpers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"lineblocs.com/api/mocks"
	"lineblocs.com/api/model"
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

		mockUserStore := mocks.UserStoreInterface{}
		mockCallStore := mocks.CallStoreInterface{}

		mockResponse := &model.Workspace{
			Id:              1,
			CreatorId:       123,
			OutboundMacroId: 123,
			Name:            "test",
		}

		mockUserStore.EXPECT().GetUserByDID("").Return("", nil)
		mockCallStore.EXPECT().GetWorkspaceByDomain("").Return(mockResponse, nil)
		mockUserStore.EXPECT().GetWorkspaceParams(1).Return(nil, nil)

		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, nil, &mockUserStore)

		if assert.NoError(t, handler.GetUserByDID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return error for UserStore", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/getUserByDID", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().GetUserByDID("").Return("", errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)

		if assert.NoError(t, handler.GetUserByDID(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return error for CallStore", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/getUserByDID", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockCallStore := mocks.CallStoreInterface{}

		mockCallStore.EXPECT().GetWorkspaceByDomain("").Return(nil, errors.New("error"))
		mockUserStore.EXPECT().GetUserByDID("").Return("", nil)
		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, nil, &mockUserStore)

		if assert.NoError(t, handler.GetUserByDID(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return error for UserStore in GetWorkspaceParams", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/user/getUserByDID", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockCallStore := mocks.CallStoreInterface{}

		mockResponse := &model.Workspace{
			Id: 1,
		}

		mockCallStore.EXPECT().GetWorkspaceByDomain("").Return(mockResponse, nil)
		mockUserStore.EXPECT().GetUserByDID("").Return("", nil)
		mockUserStore.EXPECT().GetWorkspaceParams(1).Return(nil, errors.New("error"))

		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, nil, &mockUserStore)

		if assert.NoError(t, handler.GetUserByDID(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestGetSettings(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/getSettings", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().GetSettings().Return(nil, nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.GetSettings(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/getSettings", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().GetSettings().Return(nil, errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.GetSettings(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return error for sqlNoRows", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/getSettings", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().GetSettings().Return(nil, sql.ErrNoRows)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.GetSettings(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestStoreRegistration(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return ok for no expires", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/StoreRegistration?expires=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockCallStore := mocks.CallStoreInterface{}
		mockUserStore := mocks.UserStoreInterface{}

		mockCallStore.EXPECT().GetWorkspaceByDomain("").Return(nil, nil)
		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.StoreRegistration(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return internal error for wrong domain", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/StoreRegistration?domain=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockCallStore := mocks.CallStoreInterface{}
		mockUserStore := mocks.UserStoreInterface{}

		mockCallStore.EXPECT().GetWorkspaceByDomain("abc").Return(nil, errors.New("error"))
		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.StoreRegistration(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return internal error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/StoreRegistration?expires=123&user=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockCallStore := mocks.CallStoreInterface{}
		mockUserStore := mocks.UserStoreInterface{}

		workspace := &model.Workspace{}

		mockCallStore.EXPECT().GetWorkspaceByDomain("").Return(workspace, nil)
		mockUserStore.EXPECT().StoreRegistration("abc", 123, workspace).Return(errors.New("error"))

		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.StoreRegistration(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return no errors", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/StoreRegistration?expires=123&user=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockCallStore := mocks.CallStoreInterface{}
		mockUserStore := mocks.UserStoreInterface{}

		workspace := &model.Workspace{}

		mockCallStore.EXPECT().GetWorkspaceByDomain("").Return(workspace, nil)
		mockUserStore.EXPECT().StoreRegistration("abc", 123, workspace).Return(nil)

		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.StoreRegistration(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}

func TestIncomingMediaServerValidation(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/IncomingMediaServerValidation?source=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().IncomingMediaServerValidation("abc").Return(true, nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.IncomingMediaServerValidation(c)) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}
	})

	t.Run("Should return error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/IncomingMediaServerValidation?source=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().IncomingMediaServerValidation("abc").Return(false, errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.IncomingMediaServerValidation(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return internal error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/IncomingMediaServerValidation?source=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().IncomingMediaServerValidation("abc").Return(false, nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.IncomingMediaServerValidation(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestLookupSIPTrunkByDID(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/LookupSIPTrunkByDID?did=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().LookupSIPTrunkByDID("abc").Return([]byte("abc"), nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.LookupSIPTrunkByDID(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return error for no sips", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/LookupSIPTrunkByDID?did=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().LookupSIPTrunkByDID("abc").Return(nil, nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.LookupSIPTrunkByDID(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/LookupSIPTrunkByDID?did=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().LookupSIPTrunkByDID("abc").Return(nil, errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.LookupSIPTrunkByDID(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestIncomingTrunkValidation(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")
	iptest := "0.0.0.0"

	t.Run("Should return no error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/IncomingTrunkValidation?fromdomain=0.0.0.0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().IncomingTrunkValidation(iptest).Return([]byte("abc"), nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.IncomingTrunkValidation(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/IncomingTrunkValidation?fromdomain=0.0.0.0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().IncomingTrunkValidation(iptest).Return([]byte("abc"), errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.IncomingTrunkValidation(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return error for no matches", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/IncomingTrunkValidation?fromdomain=0.0.0.0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserStore := mocks.UserStoreInterface{}
		mockUserStore.EXPECT().IncomingTrunkValidation(iptest).Return(nil, nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, &mockUserStore)
		if assert.NoError(t, handler.IncomingTrunkValidation(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}
