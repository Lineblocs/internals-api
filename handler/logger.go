package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

/*
Input: Log model
Todo : Create log model and store to db, send log email
Output: If success return NoContent else return err
*/
func (h *Handler) CreateLog(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "CreateLog is called...")
	var logReq model.Log

	if err := c.Bind(&logReq); err != nil {
		return utils.HandleInternalErr("CreateLog could not decode JSON", err, c)
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
	log := &model.LogRoutine{From: from,
		To:          to,
		Level:       level,
		Title:       logReq.Title,
		Report:      logReq.Report,
		FlowId:      logReq.FlowId,
		UserId:      logReq.UserId,
		WorkspaceId: logReq.WorkspaceId}

	workspace, err := h.callStore.GetWorkspaceFromDB(log.WorkspaceId)
	if err != nil {
		return utils.HandleInternalErr("Could not get workspace..", err, c)
	}

	_, err = h.loggerStore.StartLogRoutine(workspace, log)
	if err != nil {
		return utils.HandleInternalErr("CreateLog 2 log routine error", err, c)
	}
	return c.NoContent(http.StatusOK)
}

/*
Input: type, level, domain
Todo : Create log model and store to db, send log email
Output: If success return NoContent else return err
*/
func (h *Handler) CreateLogSimple(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "CreateLogSimple is called...")

	logType := c.FormValue("type")
	level := c.FormValue("level")
	domain := c.FormValue("domain")
	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("Could not get workspace..", err, c)
	}

	if level == "" {
		level = "infO"
	}

	var title string
	var report string
	switch logType {
	case "verify-callerid-cailed":
		title = "Caller ID Verify failed.."
		report = "Caller ID Verify failed.."
	}
	log := &model.LogRoutine{
		From:        "",
		To:          "",
		Level:       level,
		Title:       title,
		Report:      report,
		UserId:      workspace.CreatorId,
		WorkspaceId: workspace.Id}

	_, err = h.loggerStore.StartLogRoutine(workspace, log)
	if err != nil {
		return utils.HandleInternalErr("CreateLog log routine error", err, c)
	}
	return c.NoContent(http.StatusOK)
}
