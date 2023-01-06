package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

func (h *Handler) CreateFax(c echo.Context) error {
	var fax *model.Fax
	file, err := c.FormFile("file")

	if err != nil {
		return utils.HandleInternalErr("CreateFax error occured", err, c)
	}

	userId := c.FormValue("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return utils.HandleInternalErr("CreateFax error occured user ID", err, c)
	}

	workspaceId := c.FormValue("workspace_id")
	workspaceIdInt, err := strconv.Atoi(workspaceId)
	if err != nil {
		return utils.HandleInternalErr("CreateFax error occured workspace ID", err, c)
	}

	workspace, err := h.callStore.GetWorkspaceFromDB(workspaceIdInt)
	if err != nil {
		fmt.Printf("could not get workspace..")
		return c.NoContent(http.StatusNoContent)
	}

	callId := c.FormValue("call_id")
	callIdInt, err := strconv.Atoi(callId)
	if err != nil {
		return utils.HandleInternalErr("CreateFax error occured call ID", err, c)
	}

	name := c.FormValue("name")

	src, err := file.Open()
	if err != nil {
		return utils.HandleInternalErr("CreateFax error occured", err, c)
	}
	defer src.Close()

	dst, err := os.OpenFile(file.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return utils.HandleInternalErr("CreateFax error occured", err, c)
	}
	defer dst.Close()
	apiId := utils.CreateAPIID("fax")
	uri := utils.CreateS3URL("faxes", apiId)
	count, err := h.faxStore.GetFaxCount(workspaceIdInt)
	if err != nil {
		return utils.HandleInternalErr("CreateFax error occured", err, c)
	}

	fax = &model.Fax{UserId: userIdInt, WorkspaceId: workspaceIdInt, CallId: callIdInt, Uri: uri}

	faxId, err := h.faxStore.CreateFax(fax, name, file.Size, apiId, workspace.Plan)

	if err != nil {
		return utils.HandleInternalErr("CreateFax error occured", err, c)
	}

	limit, err := utils.GetPlanFaxLimit(workspace)
	if err != nil {
		return utils.HandleInternalErr("CreateFax error occured", err, c)
	}
	newCount := (*count) + 1
	if newCount > *limit {
		fmt.Printf("not saving fax due to limit reached..")
		return c.NoContent(http.StatusNoContent)
	}
	go utils.UploadS3("faxes", apiId, src)

	c.Response().Writer.Header().Set("X-Fax-ID", strconv.FormatInt(faxId, 10))
	return c.JSON(http.StatusOK, &fax)
}
