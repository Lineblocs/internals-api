package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
	"lineblocs.com/api/extension"
)

func applyExtensionValues(extensionStore extension.Store, call *model.Call, extensionValue string) (error) {
	exten, err := extensionStore.GetExtensionByUsername(call.WorkspaceId, extensionValue)
	if err != nil {
		return err
	}

	switch call.Direction {
	case "inbound":
		call.ToExtensionId = &exten.Id
	case "outbound":
		call.FromExtensionId = &exten.Id
	default:
	}
	return nil
}

/*
Input: Call model
Todo : Create new call and store to db
Output: If success return created Call model with callid in header else return err
*/
func (h *Handler) CreateCall(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "CreateCall is called...")

	var call model.Call

	if err := c.Bind(&call); err != nil {
		return utils.HandleInternalErr("CreateCall 1 Could not decode JSON", err, c)
	}
	if err := c.Validate(&call); err != nil {
		return utils.HandleInternalErr("CreateCall 1 Could not decode JSON", err, c)
	}

	call.APIId = utils.CreateAPIID("call")

	if call.Direction == "outbound" {
		// Check if this is the first time we are making a call to this destination
		go h.callStore.CheckIsMakingOutboundCallFirstTime(call)
	}

	populateExtensionDetails := c.QueryParam("populateExtensionDetails")
	extensionValue := c.QueryParam("extensionValue")
	if populateExtensionDetails != "" && populateExtensionDetails == "1" {
		err := applyExtensionValues(h.extensionStore, &call, extensionValue)
		if err != nil {
			utils.Log(logrus.InfoLevel, "CreateCall populate extension details failed. err: " + err.Error())
		}
	}
	// check if we need to populate the extension details
	callId, err := h.callStore.CreateCall(&call)
	if err != nil {
		return utils.HandleInternalErr("CreateCall Could not execute query", err, c)
	}

	c.Response().Writer.Header().Set("X-Call-ID", callId)
	return c.JSON(http.StatusOK, &call)
}

/*
Input: CallUpdate model
Todo : Update existing call with matching id
Output: If success return NoContent else return err
*/
func (h *Handler) UpdateCall(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "UpdateCall is called...")

	var update model.CallUpdate

	if err := c.Bind(&update); err != nil {
		return utils.HandleInternalErr("UpdateCall 1 Could not decode JSON", err, c)
	}
	if err := c.Validate(&update); err != nil {
		return utils.HandleInternalErr("UpdateCall 2 Could not decode JSON", err, c)
	}

	// Only update if status is "ended"
	if update.Status == "ended" {
		err := h.callStore.UpdateCall(&update)
		if err != nil {
			return utils.HandleInternalErr("UpdateCall Could not execute query..", err, c)
		}
	}

	return c.NoContent(http.StatusNoContent)
}

/*
Input: id
Todo : Fetch a call with call_id
Output: If success return Call model else return err
*/
func (h *Handler) FetchCall(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "FetchCall is called...")

	id := c.QueryParam("id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		return utils.HandleInternalErr("FetchCall error occured", err, c)
	}

	// Get call data from db with id
	call, err := h.callStore.GetCallFromDB(id_int)
	if err != nil {
		return utils.HandleInternalErr("FetchCall error occured", err, c)
	}
	return c.JSON(http.StatusOK, &call)
}

/*
Input: callid, apiid
Todo : Set sip_call_id field with matching id
Output: If success return NoContent else return err
*/
func (h *Handler) SetSIPCallID(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "SetSIPCallId is called...")

	callid := c.FormValue("callid")
	apiid := c.FormValue("apiid")

	err := h.callStore.SetSIPCallID(callid, apiid)
	if err != nil {
		return utils.HandleInternalErr("SetSIPCallID could not execute query..", err, c)
	}
	return c.NoContent(http.StatusOK)
}

/*
Input: ip, apiid
Todo : Update provider_id of call table with matching ip address
Output: If success return NoContent else return err
*/
func (h *Handler) SetProviderByIP(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "SetProviderByIP is called...")

	ip := c.FormValue("ip")
	apiid := c.FormValue("apiid")
	err := h.callStore.SetProviderByIP(ip, apiid)
	if err != nil {
		return utils.HandleInternalErr("SetProviderByID could not execute query..", err, c)
	}
	return c.NoContent(http.StatusOK)
}

/*
Input: Conference model
Todo : Create new conference and store to db
Output: If success return created Conference model with conferenceId in header else return err
*/
func (h *Handler) CreateConference(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "CreateConference is called...")

	var conference model.Conference

	if err := c.Bind(&conference); err != nil {
		return utils.HandleInternalErr("CreateConference 1 Could not decode JSON", err, c)
	}
	if err := c.Validate(&conference); err != nil {
		return utils.HandleInternalErr("CreateConference 2 Could not decode JSON", err, c)
	}

	conferenceId, err := h.callStore.CreateConference(&conference)

	if err != nil {
		return utils.HandleInternalErr("CreateConference error occured", err, c)
	}

	c.Response().Writer.Header().Set("X-Conference-ID", conferenceId)
	return c.JSON(http.StatusOK, &conference)
}
