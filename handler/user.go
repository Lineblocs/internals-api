package handler

import (
	"database/sql"
	"errors"
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

func (h *Handler) GetPSTNProviderIP(c echo.Context) error {
	fmt.Printf("received PSTN request..\r\n")
	from := c.Param("from")
	to := c.Param("to")
	domain := c.Param("domain")
	//ru := getQueryVariable(r, "ru")
	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("GetPSTNProviderIP error", err, c)
	}

	if workspace.BYOEnabled {
		info, err := h.userStore.GetBYOPSTNProvider(from, to, workspace.Id)
		if err != nil {
			return utils.HandleInternalErr("GetPSTNProviderIP error", err, c)
		}
		return c.JSON(http.StatusOK, &info)
	}

	info, err := h.userStore.GetBestPSTNProvider(from, to)
	if err != nil {
		return utils.HandleInternalErr("getPSTNProviderIp error 1 ", err, c)
	}

	return c.JSON(http.StatusOK, &info)
}

func (h *Handler) GetPSTNProviderIPForTrunk(c echo.Context) error {
	fmt.Printf("received PSTN request for trunk..\r\n")
	from := c.Param("from")
	to := c.Param("to")

	info, err := h.userStore.GetBestPSTNProvider(from, to)
	if err != nil {
		return utils.HandleInternalErr("getpstnprovideripfortrunk error 1", err, c)
	}

	return c.JSON(http.StatusOK, &info)
}

func (h *Handler) IPWhitelistLookup(c echo.Context) error {
	source := c.Param("ip")
	domain := c.Param("domain")
	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("IPWhitelistLookup error occured", err, c)
	}
	match, err := h.userStore.IPWhitelistLookup(source, workspace)
	if err != nil {
		return utils.HandleInternalErr("IPWhitelistLookup error", err, c)
	}
	if match {
		return c.NoContent(http.StatusNoContent)
	}
	return c.NoContent(http.StatusNotFound)
}

func (h *Handler) GetDIDAcceptOption(c echo.Context) error {
	did := c.Param("did")
	result, err := h.userStore.GetDIDAcceptOption(did)
	if err != nil {
		return utils.HandleInternalErr("GetDIDAcceptOption error occured", err, c)
	}
	return c.JSON(http.StatusOK, result)
}

func (h *Handler) GetDIDAssignedIP(c echo.Context) error {
	server, err := utils.GetDIDRoutedServer2(false)
	if err != nil {
		return utils.HandleInternalErr("GetUserAssignedIP error occured", err, c)
	}
	if server == nil {
		return utils.HandleInternalErr("GetUserAssignedIP could not get server", err, c)
	}
	return c.JSON(http.StatusOK, []byte(server.PrivateIpAddress))
}

func (h *Handler) GetUserAssignedIP(c echo.Context) error {
	fmt.Printf("Get assigned IP called..\r\n")
	opt := c.Param("rtcOptimized")
	var err error
	var rtcOptimized bool
	// default
	rtcOptimized = false

	if &opt != nil {
		rtcOptimized, err = strconv.ParseBool(opt)
	}
	if err != nil {
		return utils.HandleInternalErr("GetUserAssignedIP error", err, c)
	}

	domain := c.Param("domain")
	routerip := c.Param("routerip")
	fmt.Printf("Finding server for domain " + domain + "..\r\n")
	fmt.Printf("Router IP is " + routerip + "..\r\n")
	//ru := c.Param("ru")
	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("GetUserAssignedIP error 1", err, c)
	}

	server, err := h.userStore.GetUserRoutedServer2(rtcOptimized, workspace, routerip)

	if err != nil {
		return utils.HandleInternalErr("GetUserAssignedIP error occured 3", err, c)
	}
	if server == nil {
		return utils.HandleInternalErr("GetUserAssignedIP could not get server", err, c)
	}
	fmt.Printf("Found server " + server.PrivateIpAddress + "..\r\n")
	return c.JSON(http.StatusOK, []byte(server.PrivateIpAddress))
}

func (h *Handler) GetTrunkAssignedIP(c echo.Context) error {
	server, err := utils.GetDIDRoutedServer2(false)
	if err != nil {
		return utils.HandleInternalErr("GetUserAssignedIP error occured", err, c)
	}
	if server == nil {
		return utils.HandleInternalErr("GetUserAssignedIP could not get server", err, c)
	}
	return c.JSON(http.StatusOK, []byte(server.PrivateIpAddress))
}

func (h *Handler) AddPSTNProviderTechPrefix(c echo.Context) error {
	//To do
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetCallerIdToUse(c echo.Context) error {
	domain := c.Param("domain")
	extension := c.Param("extension")
	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("GetCallerIdToUse error 1 ", err, c)
	}

	callerId, err := h.userStore.GetCallerIdToUse(workspace, extension)
	if err == sql.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	}
	info := &model.CallerIDInfo{CallerID: callerId}
	return c.JSON(http.StatusOK, &info)
}

func (h *Handler) GetExtensionFlowInfo(c echo.Context) error {
	extension := c.Param("extension")
	workspaceId := c.Param("workspace")

	info, err := h.userStore.GetExtensionFlowInfo(workspaceId, extension)

	if err == sql.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, &info)
}

func (h *Handler) GetFlowInfo(c echo.Context) error {
	flowId := c.Param("flow_id")
	workspaceId := c.Param("workspace")

	info, err := h.userStore.GetFlowInfo(workspaceId, flowId)

	if err == sql.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, &info)
}

func (h *Handler) GetDIDDomain(c echo.Context) error {
	// To do
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetCodeFlowInfo(c echo.Context) error {
	code := c.Param("code")
	workspaceId := c.Param("workspace")

	info, err := h.userStore.GetCodeFlowInfo(workspaceId, code)

	if err != nil {
		return utils.HandleInternalErr("GetCodeFlowInfo error", err, c)
	}
	return c.JSON(http.StatusOK, &info)
}

func (h *Handler) IncomingDIDValidation(c echo.Context) error {
	// check for all types of call routing scenarios
	// 1. PSTN DID call
	// 2. Hosted SIP trunk call
	// 3. BYOC trunk call
	did := c.Param("did")
	number := c.Param("number")
	source := c.Param("source")

	info, err := h.userStore.IncomingDIDValidation(did)
	if err == nil {

		// check if we're routing to user SIP turnk
		if info.TrunkId != 0 {
			fmt.Printf("found trunk associated with DID number -- routing to user SIP trunk")
			return c.JSON(http.StatusOK, []byte("user_sip_trunk"))
		}
		match, err := h.userStore.CheckPSTNIPWhitelist(did, source)
		if err != nil {
			return utils.HandleInternalErr("IncomingDIDValidation error 1", err, c)
		}

		if !match {
			return utils.HandleInternalErr("IncomingDIDValidation no match found 1", err, c)
		}
		fmt.Printf("Matched incoming DID..")
		valid, err := h.userStore.FinishValidation(number, info.DidWorkspaceId)
		if err != nil {
			return utils.HandleInternalErr("IncomingDIDValidation error 2 valid", err, c)
		}
		if !valid {
			return utils.HandleInternalErr("number not valid..", err, c)
		}
		return c.JSON(http.StatusOK, []byte("network_managed"))
	}

	fmt.Println("looking up BYO DIDs now...")
	byoInfo, err := h.userStore.IncomingBYODIDValidation(did)
	if err == nil {
		match, err := h.userStore.CheckBYOPSTNIPWhitelist(did, source)
		if err != nil {
			return utils.HandleInternalErr("IncomingDIDValidation error 3", err, c)
		}
		if !match {
			return utils.HandleInternalErr("IncomingDIDValidation no match found 2", err, c)
		}
		fmt.Printf("Matched incoming DID..")
		valid, err := h.userStore.FinishValidation(number, byoInfo.DidWorkspaceId)
		if err != nil {
			return utils.HandleInternalErr("IncomingDIDValidation error 4 valid", err, c)
		}
		if !valid {
			return utils.HandleInternalErr("number not valid..", err, c)
		}

		return c.JSON(http.StatusOK, []byte("byo_carrier"))
	}
	return utils.HandleInternalErr("IncomingDIDValidation error 3", errors.New("no DID match found..."), c)
}

func (h *Handler) IncomingTrunkValidation(c echo.Context) error {
	//did := c.Param("did")
	//number := c.Param("number")
	//source := c.Param("source")
	fromdomain := c.Param("fromdomain")
	//destDomain := c.Param("destdomain")

	trunkip, err := utils.LookupSIPAddress(fromdomain)
	if err != nil {
		return utils.HandleInternalErr("IncomingTrunkValidation error 4 valid", err, c)
	}

	fmt.Printf("from domain %s trunk IP is %s..\r\n", fromdomain, *trunkip)

	result, err := h.userStore.IncomingTrunkValidation(*trunkip)
	if err != nil {
		return utils.HandleInternalErr("IncomingTrunkValidation error 1 valid", err, c)
	}

	if result == nil {
		return utils.HandleInternalErr("checked all SIP trunks no matches were found... error.", err, c)
	}
	return c.JSON(http.StatusOK, result)
}

func (h *Handler) LookupSIPTrunkByDID(c echo.Context) error {
	did := c.Param("did")

	result, err := h.userStore.LookupSIPTrunkByDID(did)
	if err != nil {
		return utils.HandleInternalErr("LookupSIPTrunkByDID error", err, c)
	}

	if result == nil {
		return utils.HandleInternalErr("checked all SIP trunks and found that none were online... error.", err, c)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) IncomingMediaServerValidation(c echo.Context) error {
	//number:= c.Param("number")
	source := c.Param("source")
	//did := c.Param("did")

	result, err := h.userStore.IncomingMediaServerValidation(source)

	if err != nil {
		return utils.HandleInternalErr("IncomingMediaServerValidation error", err, c)
	}

	if result {
		return c.NoContent(http.StatusNoContent)
	}
	return utils.HandleInternalErr("No media server found..", err, c)
}

func (h *Handler) StoreRegistration(c echo.Context) error {
	domain := c.FormValue("domain")
	//ip := rc.FormValue("ip")
	user := c.FormValue("user")
	//contact := c.FormValue("contact")
	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	var expires int

	expires, err = strconv.Atoi(c.FormValue("expires"))

	if err != nil {
		fmt.Printf("could not get expiry.. not setting online\r\n")
		return c.NoContent(http.StatusNoContent)
	}
	if err != nil {
		return utils.HandleInternalErr("StoreRegistration error..", err, c)
	}

	err = h.userStore.StoreRegistration(user, expires, workspace)
	if err != nil {
		return utils.HandleInternalErr("StoreRegistration Could not execute query..", err, c)
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetSettings(c echo.Context) error {
	settings, err := h.userStore.GetSettings()
	if err == sql.ErrNoRows {
		// no records setup were setup, just return empty
		return utils.HandleInternalErr("GetSettings no rows found..", err, c)
	}
	if err != nil {
		return utils.HandleInternalErr("GetSettings error:"+err.Error(), err, c)
	}
	return c.JSON(http.StatusOK, &settings)
}

func (h *Handler) ProcessSIPTrunkCall(c echo.Context) error {
	did := c.Param("did")

	result, err := h.userStore.ProcessSIPTrunkCall(did)
	if err != nil {
		return utils.HandleInternalErr("ProcessSIPTrunkCall error 1 valid", err, c)
	}

	if result != nil {
		return c.JSON(http.StatusOK, &result)
	}

	return utils.HandleInternalErr("No trunks to route to..", err, c)
}
