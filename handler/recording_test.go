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

func TestCreateRecording(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no error", func(t *testing.T) {

		recBody := model.Recording{
			Id:          1,
			WorkspaceId: 1,
			APIId:       "rec-7700845c-72b7-415b-a8ee-996b8d4ec239",
		}

		bodyBytes, _ := json.Marshal(recBody)

		req := httptest.NewRequest(http.MethodGet, "/user/CreateRecording", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockWorkspace := model.Workspace{}

		mockCallStore := mocks.CallStoreInterface{}
		mockRecStore := mocks.RecordingStoreInterface{}
		mockCallStore.EXPECT().GetWorkspaceFromDB(1).Return(&mockWorkspace, nil)
		mockRecStore.EXPECT().CreateRecording(&mockWorkspace, mock.Anything).Return(1, nil)

		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, &mockRecStore, nil)
		if assert.NoError(t, handler.CreateRecording(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return bind error", func(t *testing.T) {

		type InvalidStruct struct {
			Id string `json:"id"`
		}

		recBody := InvalidStruct{
			Id: "abc",
		}

		bodyBytes, _ := json.Marshal(recBody)

		req := httptest.NewRequest(http.MethodGet, "/user/CreateRecording", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := NewHandler(nil, nil, nil, nil, nil, nil, nil, nil)
		if assert.NoError(t, handler.CreateRecording(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return workspace error", func(t *testing.T) {

		recBody := model.Recording{
			Id:          1,
			WorkspaceId: 1,
			APIId:       "rec-7700845c-72b7-415b-a8ee-996b8d4ec239",
		}

		bodyBytes, _ := json.Marshal(recBody)

		req := httptest.NewRequest(http.MethodGet, "/user/CreateRecording", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockWorkspace := model.Workspace{}

		mockCallStore := mocks.CallStoreInterface{}
		mockCallStore.EXPECT().GetWorkspaceFromDB(1).Return(&mockWorkspace, errors.New("error"))

		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, nil, nil)
		if assert.NoError(t, handler.CreateRecording(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return recording error", func(t *testing.T) {

		recBody := model.Recording{
			Id:          1,
			WorkspaceId: 1,
			APIId:       "rec-7700845c-72b7-415b-a8ee-996b8d4ec239",
		}

		bodyBytes, _ := json.Marshal(recBody)

		req := httptest.NewRequest(http.MethodGet, "/user/CreateRecording", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockWorkspace := model.Workspace{}

		mockCallStore := mocks.CallStoreInterface{}
		mockRecStore := mocks.RecordingStoreInterface{}
		mockCallStore.EXPECT().GetWorkspaceFromDB(1).Return(&mockWorkspace, nil)
		mockRecStore.EXPECT().CreateRecording(&mockWorkspace, mock.Anything).Return(1, errors.New("errors"))

		handler := NewHandler(nil, &mockCallStore, nil, nil, nil, nil, &mockRecStore, nil)
		if assert.NoError(t, handler.CreateRecording(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestGetRecording(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/GetRecording?id=123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRecStore := mocks.RecordingStoreInterface{}
		mockRecStore.EXPECT().GetRecordingFromDB(123).Return(nil, nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, &mockRecStore, nil)
		if assert.NoError(t, handler.GetRecording(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/GetRecording?id=abc", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRecStore := mocks.RecordingStoreInterface{}
		handler := NewHandler(nil, nil, nil, nil, nil, nil, &mockRecStore, nil)
		if assert.NoError(t, handler.GetRecording(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return recording error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/GetRecording?id=123", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRecStore := mocks.RecordingStoreInterface{}
		mockRecStore.EXPECT().GetRecordingFromDB(123).Return(nil, errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, &mockRecStore, nil)
		if assert.NoError(t, handler.GetRecording(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestUpdateRecordingTranscription(t *testing.T) {

	e := echo.New()
	helpers.InitLogrus("stdout")

	t.Run("Should return no content", func(t *testing.T) {

		recBody := model.RecordingTranscription{
			RecordingId: 1,
		}

		bodyBytes, _ := json.Marshal(recBody)
		req := httptest.NewRequest(http.MethodGet, "/user/updateRecordingTranscription", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRecStore := mocks.RecordingStoreInterface{}
		mockRecStore.EXPECT().UpdateRecordingTranscription(&recBody).Return(nil)
		handler := NewHandler(nil, nil, nil, nil, nil, nil, &mockRecStore, nil)
		if assert.NoError(t, handler.UpdateRecordingTranscription(c)) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}
	})

	t.Run("Should return upload error", func(t *testing.T) {

		recBody := model.RecordingTranscription{
			RecordingId: 1,
		}

		bodyBytes, _ := json.Marshal(recBody)
		req := httptest.NewRequest(http.MethodGet, "/user/UpdateRecordingTranscription", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRecStore := mocks.RecordingStoreInterface{}
		mockRecStore.EXPECT().UpdateRecordingTranscription(&recBody).Return(errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, &mockRecStore, nil)
		if assert.NoError(t, handler.UpdateRecordingTranscription(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Should return bind error", func(t *testing.T) {

		type InvalidStruct struct {
			RecordingId string `json:"recording_id"`
		}

		recBody := InvalidStruct{
			RecordingId: "abc",
		}

		bodyBytes, _ := json.Marshal(recBody)
		req := httptest.NewRequest(http.MethodGet, "/user/UpdateRecordingTranscription", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRecStore := mocks.RecordingStoreInterface{}
		mockRecStore.EXPECT().UpdateRecordingTranscription(nil).Return(errors.New("error"))
		handler := NewHandler(nil, nil, nil, nil, nil, nil, &mockRecStore, nil)

		if assert.NoError(t, handler.UpdateRecordingTranscription(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}
