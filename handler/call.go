package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

func (h *Handler) CreateCall(c echo.Context) error {
	var call model.Call

	if err := c.Bind(&call); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := c.Validate(&call); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	call.APIId = utils.CreateAPIID("call")

	if call.Direction == "outbound" {
		//check if this is the first time we are making a call to this destination
		go h.callStore.CheckIsMakingOutboundCallFirstTime(call)
	}

	callId, err := h.callStore.CreateCall(&call)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	c.Response().Writer.Header().Set("X-Call-ID", callId)
	return c.JSON(http.StatusOK, &call)
}

func (h *Handler) UpdateCall(c echo.Context) error {
	var update model.CallUpdate

	if err := c.Bind(&update); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := c.Validate(&update); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if update.Status == "ended" {
		err := h.callStore.UpdateCall(&update)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) FetchCall(c echo.Context) error {
	id := c.Param("id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	call, err := h.callStore.GetCallFromDB(id_int)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, &call)
}

func (h *Handler) SetSIPCallID(c echo.Context) error {
	callid := c.FormValue("callid")
	apiid := c.FormValue("apiid")

	err := h.callStore.SetSIPCallID(callid, apiid)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) SetProviderByIP(c echo.Context) error {
	ip := c.FormValue("ip")
	apiid := c.FormValue("apiid")
	err := h.callStore.SetProviderByIP(ip, apiid)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) CreateConference(c echo.Context) error {
	var conference model.Conference

	if err := c.Bind(&conference); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := c.Validate(&conference); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	conferenceId, err := h.callStore.CreateConference(&conference)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	c.Response().Writer.Header().Set("X-Conference-ID", conferenceId)
	return c.JSON(http.StatusOK, &conference)
}
