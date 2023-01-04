package handler

import (
	"github.com/labstack/echo/v4"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

type callCreateRequest struct {
	Call struct {
		From         string `json:"from"`
		To           string `json:"to"`
		Status       string `json:"status"`
		Direction    string `json:"direction"`
		Duration     string `json:"duration"`
		UserId       int    `json:"user_id"`
		WorkspaceId  int    `json:"workspace_id"`
		APIId        string `json:"api_id"`
		SourceIp     string `json:"source_ip"`
		ChannelId    string `json:"channel_id"`
		SIPCallId    string `json:"call_id"`
		StartedAt    string `json:"started_at"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		PlanSnapshot string `json:"plan_snapshot"`
	} `json:"call"`
}

func (r *callCreateRequest) bind(c echo.Context, call *model.Call) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	call.From = r.Call.From
	call.To = r.Call.To
	call.Status = r.Call.Status
	call.Direction = r.Call.Direction
	call.UserId = r.Call.UserId
	call.WorkspaceId = r.Call.WorkspaceId
	call.APIId = utils.CreateAPIID("call")
	call.ChannelId = r.Call.ChannelId
	call.SIPCallId = r.Call.SIPCallId
	return nil
}
