package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	helpers "github.com/Lineblocs/go-helpers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"lineblocs.com/api/mocks"
	"lineblocs.com/api/model"
)

func TestCreateLogSimple(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return domain error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/CreateLogSimple", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockCallStore := mocks.CallStoreInterface{}
		mockCallStore.EXPECT().GetWorkspaceByDomain("").Return(nil, errors.New("error"))
		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, nil, nil)
		if assert.NoError(t, handler.CreateLogSimple(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return no error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/CreateLogSimple", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockResponse := &model.Workspace{
			Id:        1,
			CreatorId: 1,
		}

		mockCallStore := mocks.CallStoreInterface{}
		mockLoggerStore := mocks.LoggerStoreInterface{}
		mockCallStore.EXPECT().GetWorkspaceByDomain("").Return(mockResponse, nil)
		mockLoggerStore.EXPECT().StartLogRoutine(mock.Anything, mock.Anything).Return(nil, nil)
		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, &mockLoggerStore, nil, nil)
		if assert.NoError(t, handler.CreateLogSimple(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return log routine error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/createLogSimple?type=verify-callerid-cailed", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockResponse := &model.Workspace{
			Id:        1,
			CreatorId: 1,
		}

		mockCallStore := mocks.CallStoreInterface{}
		mockLoggerStore := mocks.LoggerStoreInterface{}
		mockCallStore.EXPECT().GetWorkspaceByDomain("").Return(mockResponse, nil)
		mockLoggerStore.EXPECT().StartLogRoutine(mock.Anything, mock.Anything).Return(nil, errors.New("error"))
		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, &mockLoggerStore, nil, nil)
		if assert.NoError(t, handler.CreateLogSimple(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestCreateLog(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return bind error", func(t *testing.T) {

		type InvalidStruct struct {
			UserId string `json:"user_id"`
		}

		recBody := InvalidStruct{
			UserId: "abc",
		}

		bodyBytes, _ := json.Marshal(recBody)

		req := httptest.NewRequest(http.MethodGet, "/user/CreateLog", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockCallStore := mocks.CallStoreInterface{}
		mockCallStore.EXPECT().GetWorkspaceByDomain("").Return(nil, errors.New("error"))
		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, nil, nil)
		if assert.NoError(t, handler.CreateLog(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return no error", func(t *testing.T) {

		info := "info"
		from := "from"
		to := "to"

		recBody := model.Log{
			UserId:      1,
			WorkspaceId: 1,
			Title:       "",
			Report:      "",
			FlowId:      0,
			Level:       &info,
			From:        &from,
			To:          &to,
		}

		bodyBytes, _ := json.Marshal(recBody)

		req := httptest.NewRequest(http.MethodGet, "/user/CreateLog", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockResponse := &model.Workspace{
			Id:        1,
			CreatorId: 1,
		}

		mockCallStore := mocks.CallStoreInterface{}
		mockLoggerStore := mocks.LoggerStoreInterface{}
		mockCallStore.EXPECT().GetWorkspaceFromDB(1).Return(mockResponse, nil)
		mockLoggerStore.EXPECT().StartLogRoutine(mock.Anything, mock.Anything).Return(nil, nil)
		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, &mockLoggerStore, nil, nil)
		if assert.NoError(t, handler.CreateLog(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return workspace error", func(t *testing.T) {

		info := "info"
		from := "from"
		to := "to"

		recBody := model.Log{
			UserId:      1,
			WorkspaceId: 1,
			Title:       "",
			Report:      "",
			FlowId:      0,
			Level:       &info,
			From:        &from,
			To:          &to,
		}

		bodyBytes, _ := json.Marshal(recBody)

		req := httptest.NewRequest(http.MethodGet, "/user/CreateLog", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockCallStore := mocks.CallStoreInterface{}
		mockLoggerStore := mocks.LoggerStoreInterface{}
		mockCallStore.EXPECT().GetWorkspaceFromDB(1).Return(nil, errors.New("error"))
		mockLoggerStore.EXPECT().StartLogRoutine(mock.Anything, mock.Anything).Return(nil, nil)
		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, &mockLoggerStore, nil, nil)
		if assert.NoError(t, handler.CreateLog(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return log routine error", func(t *testing.T) {

		info := "info"
		from := "from"
		to := "to"

		recBody := model.Log{
			UserId:      1,
			WorkspaceId: 1,
			Title:       "",
			Report:      "",
			FlowId:      0,
			Level:       &info,
			From:        &from,
			To:          &to,
		}

		bodyBytes, _ := json.Marshal(recBody)

		req := httptest.NewRequest(http.MethodGet, "/user/CreateLog", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockCallStore := mocks.CallStoreInterface{}
		mockLoggerStore := mocks.LoggerStoreInterface{}
		mockCallStore.EXPECT().GetWorkspaceFromDB(1).Return(nil, nil)
		mockLoggerStore.EXPECT().StartLogRoutine(mock.Anything, mock.Anything).Return(nil, errors.New("error"))
		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, &mockLoggerStore, nil, nil)
		if assert.NoError(t, handler.CreateLog(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}
