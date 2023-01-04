package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

func (h *Handler) CreateCall(c echo.Context) error {
	var call model.Call

	req := &callCreateRequest{}
	if err := req.bind(c, &call); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if call.Direction == "outbound" {
		//check if this is the first time we are making a call to this destination
		go h.callStore.CheckIsMakingOutboundCallFirstTime(call)
	}

	callId, err := h.callStore.CreateCall(&call)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	c.Response().Writer.Header().Set("X-Call-ID", callId)
	return c.JSON(http.StatusOK, call)
}
