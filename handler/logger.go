package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

func (h *Handler) CreateLog(c echo.Context) error {
	var logReq model.Log

	if err := c.Bind(&logReq); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := c.Validate(&logReq); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	level := "info"
	if logReq.Level != nil {
		level = *logReq.Level
	}
	from := ""
	if logReq.From != nil {
		from = *logReq.From
	}

	to := ""
	if logReq.To != nil {
		to = *logReq.To
	}
	var log *model.LogRoutine = &model.LogRoutine{From: from,
		To:          to,
		Level:       level,
		Title:       logReq.Title,
		Report:      logReq.Report,
		FlowId:      logReq.FlowId,
		UserId:      logReq.UserId,
		WorkspaceId: logReq.WorkspaceId}

	workspace, err := h.callStore.GetWorkspaceFromDB(log.WorkspaceId)
	if err != nil {
		return utils.HandleInternalErr("could not get workspace..", err, c)
	}

	_, err = h.loggerStore.StartLogRoutine(workspace, log)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) CreateLogSimple(c echo.Context) error {
	logType := c.FormValue("type")
	level := c.FormValue("level")
	domain := c.FormValue("domain")
	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if &level == nil {
		level = "infO"
	}

	var title string
	var report string
	switch logType {
	case "verify-callerid-cailed":
		title = "Caller ID Verify failed.."
		report = "Caller ID Verify failed.."
	}
	var log *model.LogRoutine = &model.LogRoutine{
		From:        "",
		To:          "",
		Level:       level,
		Title:       title,
		Report:      report,
		UserId:      workspace.CreatorId,
		WorkspaceId: workspace.Id}

	_, err = h.loggerStore.StartLogRoutine(workspace, log)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.NoContent(http.StatusNoContent)
}
