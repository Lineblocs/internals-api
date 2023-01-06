package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

func (h *Handler) VerifyCaller(c echo.Context) error {
	workspaceId := c.Param("workspace_id")
	workspaceIdInt, err := strconv.Atoi(workspaceId)
	if err != nil {
		return utils.HandleInternalErr("VerifyCaller error occured", err, c)
	}
	number := c.Param("number")

	var workspace *model.Workspace

	workspace, err = h.callStore.GetWorkspaceFromDB(workspaceIdInt)
	if err != nil {
		return utils.HandleInternalErr("Workspace error occured", err, c)
	}

	valid, err := h.userStore.DoVerifyCaller(workspace, number)

	if err != nil {
		return utils.HandleInternalErr("VerifyCaller error occured", err, c)
	}
	result := model.VerifyNumber{Valid: valid}
	return c.JSON(http.StatusOK, &result)
}

func (h *Handler) VerifyCallerByDomain(c echo.Context) error {
	domain := c.Param("domain")
	number := c.Param("number")

	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("VerifyCallerByDomain error 1 occured", err, c)
	}
	valid, err := h.userStore.DoVerifyCaller(workspace, number)
	if err != nil {
		return utils.HandleInternalErr("VerifyCaller error 2 occured", err, c)
	}
	if !valid {
		return utils.HandleInternalErr("VerifyCaller number not valid", err, c)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetUserByDomain(c echo.Context) error {
	domain := c.Param("domain")

	// info, err := h.userStore.GetUserByDomain(domain)

	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("GetUserByDomain error occured", err, c)
	}

	params, err := h.userStore.GetWorkspaceParams(workspace.Id)
	if err != nil {
		return utils.HandleInternalErr("GetUserByDomain error occured", err, c)
	}

	full := &model.WorkspaceCreatorFullInfo{
		Id:              workspace.CreatorId,
		Workspace:       workspace,
		WorkspaceParams: params,
		WorkspaceName:   workspace.Name,
		WorkspaceDomain: fmt.Sprintf("%s.lineblocs.com", workspace.Name),
		WorkspaceId:     workspace.Id,
		OutboundMacroId: workspace.OutboundMacroId}

	return c.JSON(http.StatusOK, &full)
}

func (h *Handler) GetUserByDID(c echo.Context) error {
	did := c.Param("did")

	domain, err := h.userStore.GetUserByDID(did)
	if err != nil {
		return utils.HandleInternalErr("GetUserByDID error occured", err, c)
	}

	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("GetUserByDID error occured", err, c)
	}

	// Execute the query
	params, err := h.userStore.GetWorkspaceParams(workspace.Id)
	if err != nil {
		return utils.HandleInternalErr("GetUserByDID error occured", err, c)
	}
	full := &model.WorkspaceCreatorFullInfo{
		Id:              workspace.CreatorId,
		Workspace:       workspace,
		WorkspaceParams: params,
		WorkspaceName:   workspace.Name,
		WorkspaceDomain: fmt.Sprintf("%s.lineblocs.com", workspace.Name),
		WorkspaceId:     workspace.Id,
		OutboundMacroId: workspace.OutboundMacroId}

	return c.JSON(http.StatusOK, &full)
}

func (h *Handler) GetUserByTrunkSourceIp(c echo.Context) error {
	sourceIp := c.Param("source_ip")

	domain, err := h.userStore.GetUserByTrunkSourceIp(sourceIp)
	if err != nil {
		return utils.HandleInternalErr("GetUserByTrunkSourceIp error occured", err, c)
	}

	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("GetUserByTrunkSourceIp error occured", err, c)
	}

	// Execute the query
	params, err := h.userStore.GetWorkspaceParams(workspace.Id)
	if err != nil {
		return utils.HandleInternalErr("GetUserByTrunkSourceIp error occured", err, c)
	}
	full := &model.WorkspaceCreatorFullInfo{
		Id:              workspace.CreatorId,
		Workspace:       workspace,
		WorkspaceParams: params,
		WorkspaceName:   workspace.Name,
		WorkspaceDomain: fmt.Sprintf("%s.lineblocs.com", workspace.Name),
		WorkspaceId:     workspace.Id,
		OutboundMacroId: workspace.OutboundMacroId}

	return c.JSON(http.StatusOK, &full)
}

func (h *Handler) GetWorkspaceMacros(c echo.Context) error {
	workspace := c.Param("workspace")
	values, err := h.userStore.GetWorkspaceMacros(workspace)

	if err != nil {
		return utils.HandleInternalErr("GetWorkspaceMacros error", err, c)
	}
	return c.JSON(http.StatusOK, &values)
}

func (h *Handler) GetDIDNumberData(c echo.Context) error {
	number := c.Param("number")
	info, flowJson, err := h.userStore.GetDIDNumberData(number)
	if err != nil && err != sql.ErrNoRows {
		return utils.HandleInternalErr("GetDIDNumberData lookup error", err, c)
	}
	if err == sql.ErrNoRows {
		info, flowJson, err := h.userStore.GetBYODIDNumberData(number)
		if err != nil {
			return utils.HandleInternalErr("GetDIDNumberData 3 error", err, c)
		}

		if flowJson.Valid {
			info.FlowJSON = flowJson.String
		}

		params, err := h.userStore.GetWorkspaceParams(info.WorkspaceId)
		info.WorkspaceParams = params
	}
	if flowJson.Valid {
		info.FlowJSON = flowJson.String
	}

	params, err := h.userStore.GetWorkspaceParams(info.WorkspaceId)
	if err != nil {
		return utils.HandleInternalErr("GetDIDNumberData 1 error", err, c)
	}

	info.WorkspaceParams = params
	return c.JSON(http.StatusOK, &info)
}
