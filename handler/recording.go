package handler

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

func (h *Handler) CreateRecording(c echo.Context) error {
	var recording model.Recording

	if err := c.Bind(&recording); err != nil {
		return utils.HandleInternalErr("CreateCall Could not decode JSON", err, c)
	}
	if err := c.Validate(&recording); err != nil {
		return utils.HandleInternalErr("CreateCall Could not decode JSON", err, c)
	}

	recording.APIId = utils.CreateAPIID("rec")

	workspace, err := h.callStore.GetWorkspaceFromDB(recording.WorkspaceId)
	if err != nil {
		return utils.HandleInternalErr("could not get workspace..", err, c)
	}

	recId, err := h.recordingStore.CreateRecording(workspace, &recording)
	if err != nil {
		return utils.HandleInternalErr("CreateRecording error.", err, c)
	}
	c.Response().Writer.Header().Set("X-Recording-ID", strconv.FormatInt(recId, 10))
	return c.JSON(http.StatusOK, &recording)
}

func (h *Handler) UpdateRecording(c echo.Context) error {
	file, err := c.FormFile("file")
	status := c.FormValue("status")
	recordingId := c.FormValue("recording_id")
	recordingIdInt, err := strconv.Atoi(recordingId)
	record, err := h.recordingStore.GetRecordingFromDB(recordingIdInt)
	if err != nil {
		return utils.HandleInternalErr("could not get recording..", err, c)
	}

	workspace, err := h.callStore.GetWorkspaceFromDB(record.WorkspaceId)
	if err != nil {
		return utils.HandleInternalErr("could not get workspace..", err, c)
	}

	src, err := file.Open()
	if err != nil {
		return utils.HandleInternalErr("UpdateRecording error occured", err, c)
	}
	defer src.Close()

	dst, err := os.OpenFile(file.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return utils.HandleInternalErr("UpdateRecording error occured", err, c)
	}
	defer dst.Close()
	size, err := h.recordingStore.GetRecordingSpace(workspace.Id)
	if err != nil {
		return utils.HandleInternalErr("could not get recording space..", err, c)
	}
	apiId := utils.CreateAPIID("rec")
	err = h.recordingStore.UpdateRecording(apiId, status, file.Size, recordingIdInt)
	if err != nil {
		return utils.HandleInternalErr("UpdateRecording error occured", err, c)
	}

	//will not save
	limit, err := utils.GetPlanRecordingLimit(workspace)
	newSpace := size + int(file.Size)
	if newSpace > limit {
		return utils.HandleInternalErr("not saving recording due to space limit reached..", err, c)
	}
	go utils.UploadS3("recordings", apiId, src)
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) UpdateRecordingTranscription(c echo.Context) error {
	var update model.RecordingTranscription
	if err := c.Bind(&update); err != nil {
		return utils.HandleInternalErr("UpdateRecordingTranscription error", err, c)
	}
	err := h.recordingStore.UpdateRecordingTranscription(&update)
	if err != nil {
		return utils.HandleInternalErr("UpdateRecording Could not execute query", err, c)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetRecording(c echo.Context) error {
	id := c.Param("id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		return utils.HandleInternalErr("GetRecording error occured", err, c)
	}
	record, err := h.recordingStore.GetRecordingFromDB(id_int)
	if err != nil {
		return utils.HandleInternalErr("GetRecording error occured", err, c)
	}
	return c.JSON(http.StatusOK, &record)
}