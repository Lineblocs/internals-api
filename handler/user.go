package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

/*
Input: workspace_id, number
Todo : Check number is valid with domain?
Output: If success return VerifyNumber model else return err
*/
func (h *Handler) VerifyCaller(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "VerifyCaller is called...")

	workspaceId := c.QueryParam("workspace_id")
	workspaceIdInt, err := strconv.Atoi(workspaceId)
	if err != nil {
		return utils.HandleInternalErr("VerifyCaller error occured", err, c)
	}
	number := c.QueryParam("number")

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

/*
Input: domain, number
Todo : Check number is valid with domain?
Output: If success return NoContent else return err
*/
func (h *Handler) VerifyCallerByDomain(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "VerifyCallerByDomain is called...")

	domain := c.QueryParam("domain")
	number := c.QueryParam("number")

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

/*
Input: domain
Todo : Get WorkspaceCreator with matching domain and workspace
Output: If success return WorkspaceCreatorFullInfo model else return err
*/
func (h *Handler) GetUserByDomain(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetUserByDomain is called...")

	domain := c.QueryParam("domain")

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

/*
Input: did
Todo : Get WorkspaceCreator with matching DID
Output: If success return WorkspaceCreatorFullInfo model else return err
*/
func (h *Handler) GetUserByDID(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetUserByDID is called...")

	did := c.QueryParam("did")

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

/*
Input: source_ip
Todo : Get WorkspaceCreator with matching source ip
Output: If success return WorkspaceCreatorFullInfo model else return err
*/
func (h *Handler) GetUserByTrunkSourceIp(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetUserByTrunkSourceIp is called...")

	sourceIp := c.QueryParam("source_ip")

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

/*
Input: workspace
Todo : Get macro_functions with matching workspace_id
Output: If success return MacroFunction model else return err
*/
func (h *Handler) GetWorkspaceMacros(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetWorkspaceMacros is called...")

	workspace := c.QueryParam("workspace")
	values, err := h.userStore.GetWorkspaceMacros(workspace)

	if err != nil {
		return utils.HandleInternalErr("GetWorkspaceMacros error", err, c)
	}
	return c.JSON(http.StatusOK, &values)
}

/*
Input: number
Todo : Get WorkspaceDidInfo with matching number (Check both DIDNumber and BYODIDNumber)
Output: If success return WorkspaceDidInfo model else return err
*/
func (h *Handler) GetDIDNumberData(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetDIDNumberData is called...")

	number := c.QueryParam("number")
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
		return utils.HandleInternalErr("GetDIDNumberData error", err, c)
	}

	info.WorkspaceParams = params
	return c.JSON(http.StatusOK, &info)
}

/*
Input: from, to, domain
Todo : Get PSTNInfo with matching from, to, domain params
Output: If success return PSTNInfo model else return err
*/
func (h *Handler) GetPSTNProviderIP(c echo.Context) error {
	utils.Log(logrus.InfoLevel, fmt.Sprintf("Received PSTN request..\r\n"))

	from := c.QueryParam("from")
	to := c.QueryParam("to")
	domain := c.QueryParam("domain")
	//ru := c.QueryParam("ru")
	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("GetPSTNProviderIP error", err, c)
	}

	// If BYOEnabled GetBYOPSTNProvider else BestPSTNProvider
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

/*
Input: from, to
Todo : Get PSTNInfo with matching from, to params
Output: If success return PSTNInfo model else return err
*/
func (h *Handler) GetPSTNProviderIPForTrunk(c echo.Context) error {
	utils.Log(logrus.InfoLevel, fmt.Sprintf("Received PSTN request for trunk..\r\n"))
	from := c.QueryParam("from")
	to := c.QueryParam("to")

	info, err := h.userStore.GetBestPSTNProvider(from, to)
	if err != nil {
		return utils.HandleInternalErr("GetPSTNProviderIPForTrunk error", err, c)
	}

	return c.JSON(http.StatusOK, &info)
}

/*
Input: ip, domain
Todo : Check ip_whitelist with matching domain and ip
Output: If matched return StatusNoContent, not matched return StatusNotFound, error return err
*/
func (h *Handler) IPWhitelistLookup(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "IPWhitelistLookup is called...")

	source := c.QueryParam("ip")
	domain := c.QueryParam("domain")
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

/*
Input: ip, fromdomain, todomain, fromuser, touser
Todo : Check if this is a valid hosted SIP trunk or not
Output: If matched return StatusNoContent, not matched return StatusNotFound, error return err
*/
func (h *Handler) HostedSIPTrunkLookup(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "HostedSIPTrunkLookup is called...")

	source := c.QueryParam("ip")
	domain := c.QueryParam("domain")
	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("HostedSIPTrunkLookup error occured", err, c)
	}
	match, err := h.userStore.HostedSIPTrunkLookup(source, workspace)
	if err != nil {
		return utils.HandleInternalErr("HostedSIPTrunkLookup error", err, c)
	}
	if match {
		return c.NoContent(http.StatusNoContent)
	}
	return c.NoContent(http.StatusNotFound)
}

/*
Input: did
Todo : Get did_action from did_numbers or byo_did_numbers with matching did
Output: If success return did_action else return err
*/
func (h *Handler) GetDIDAcceptOption(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetDIDAcceptOption is called...")

	did := c.QueryParam("did")
	result, err := h.userStore.GetDIDAcceptOption(did)
	if err != nil {
		return utils.HandleInternalErr("GetDIDAcceptOption error occured", err, c)
	}
	return c.JSONBlob(http.StatusOK, result)
}

/*
Input:
Todo : Get DIDAssignedIP
Output: If success return PrivateIpAddress else return err
*/
func (h *Handler) GetDIDAssignedIP(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetDIDAssignedIP is called...")

	server, err := utils.GetDIDRoutedServer2(false)
	if err != nil {
		return utils.HandleInternalErr("GetUserAssignedIP error occured", err, c)
	}
	if server == nil {
		return utils.HandleInternalErr("GetUserAssignedIP could not get server", err, c)
	}
	canPlace, err := utils.CanPlaceAdditionalCalls()
	if err != nil {
		return utils.HandleInternalErr("GetUserAssignedIP error occured", err, c)
	}
	if !canPlace {
		return c.NoContent(http.StatusConflict)
	}
	return c.JSONBlob(http.StatusOK, []byte(server.PrivateIpAddress))
}

/*
Input: rtcOptimized, domain, routerip
Todo : Get UserAssignedIP
Output: If success return PrivateIpAddress else return err
*/
func (h *Handler) GetUserAssignedIP(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "Get assigned IP called..\r\n")

	opt := c.QueryParam("rtcOptimized")
	var err error
	var rtcOptimized bool
	// default
	rtcOptimized = false

	if &opt != nil && opt != "" {
		rtcOptimized, err = strconv.ParseBool(opt)
	}
	if err != nil {
		return utils.HandleInternalErr("GetUserAssignedIP error", err, c)
	}

	domain := c.QueryParam("domain")
	routerip := c.QueryParam("routerip")
	utils.Log(logrus.InfoLevel, "Finding server for domain "+domain+"..\r\n")
	utils.Log(logrus.InfoLevel, "Router IP is "+routerip+"..\r\n")
	//ru := c.QueryParam("ru")
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
	utils.Log(logrus.InfoLevel, "Found server "+server.PrivateIpAddress+"..\r\n")
	return c.JSONBlob(http.StatusOK, []byte(server.PrivateIpAddress))
}

/*
Input: rtcOptimized, domain, routerip12
Todo : Get TrunkAssignedIP
Output: If success return PrivateIpAddress else return err
*/
func (h *Handler) GetTrunkAssignedIP(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetTrunkAssignedIP is called...")

	server, err := utils.GetDIDRoutedServer2(false)
	if err != nil {
		return utils.HandleInternalErr("GetUserAssignedIP error occured", err, c)
	}
	if server == nil {
		return utils.HandleInternalErr("GetUserAssignedIP could not get server", err, c)
	}
	return c.JSONBlob(http.StatusOK, []byte(server.PrivateIpAddress))
}

func (h *Handler) AddPSTNProviderTechPrefix(c echo.Context) error {
	//To do
	return c.NoContent(http.StatusOK)
}

/*
Input: domain, extension
Todo : Get CallerId with mathcing domain and extension
Output: If successfuly find callerid return CallerIDInfo model
else return StatusNotFound(it doesn't occur error) or err(it occurs error)
*/
func (h *Handler) GetCallerIdToUse(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetCallerIdToUse is called...")

	domain := c.QueryParam("domain")
	extension := c.QueryParam("extension")
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

/*
Input: extension, workspace_id
Todo : Get ExtensionFlowInfo with matching workspace and extension
Output: If success return ExtensionFlowInfo model else return StatusNoFound or err
*/
func (h *Handler) GetExtensionFlowInfo(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetExtensionFlowInfo is called...")

	extension := c.QueryParam("extension")
	workspaceId := c.QueryParam("workspace")

	info, err := h.userStore.GetExtensionFlowInfo(workspaceId, extension)

	if err == sql.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, &info)
}

/*
Input: flow_id, workspace_id
Todo : Get ExtensionFlowInfo with matching flow_id and workspace_id
Output: If success return ExtensionFlowInfo model else return StatusNoFound or err
*/
func (h *Handler) GetFlowInfo(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetFlowInfo is called...")

	flowId := c.QueryParam("flow_id")
	workspaceId := c.QueryParam("workspace")

	info, err := h.userStore.GetFlowInfo(workspaceId, flowId)

	if err == sql.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	}

	if err != nil {
		return utils.HandleInternalErr("GetFlowInfo error", err, c)
	}

	return c.JSON(http.StatusOK, &info)
}

func (h *Handler) GetDIDDomain(c echo.Context) error {
	// To do
	return c.NoContent(http.StatusOK)
}

/*
Input: code, workspace_id
Todo : Get CodeFlowInfo with matching code and workspace_id
Output: If success return CodeFlowInfo model else return err
*/
func (h *Handler) GetCodeFlowInfo(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetCodeFlowInfo is called...")

	code := c.QueryParam("code")
	workspaceId := c.QueryParam("workspace")

	info, err := h.userStore.GetCodeFlowInfo(workspaceId, code)

	if err != nil {
		return utils.HandleInternalErr("GetCodeFlowInfo error", err, c)
	}
	return c.JSON(http.StatusOK, &info)
}

/*
Input: did, number, source
Todo : Check for all types of call routing scenarios(1.PSTN DID call, 2.Hosted SIP trunk call, 3.BYOC trunk call )
Output: If success return network_managed or byo_carrier else return err
*/
func (h *Handler) IncomingDIDValidation(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "IncomingDIDValidation is called...")

	did := c.QueryParam("did")
	number := c.QueryParam("number")
	source := c.QueryParam("source")

	info, err := h.userStore.IncomingDIDValidation(did)
	if err == nil {

		utils.Log(logrus.InfoLevel, "found DID number in database")

		// check if we're routing to user SIP turnk
		if info.TrunkId != 0 {
			utils.Log(logrus.InfoLevel, "found trunk associated with DID number -- routing to user SIP trunk")
			return c.String(http.StatusOK, "user_sip_trunk")
		}

		utils.Log(logrus.InfoLevel, "checking if IP: " + source + " is whitelisted.")
		match, err := h.userStore.CheckPSTNIPWhitelist(did, source)
		if err != nil {
			return utils.HandleInternalErr("IncomingDIDValidation error 1", err, c)
		}

		if !match {
			return utils.HandleInternalErr("IncomingDIDValidation no match found 1", err, c)
		}

		utils.Log(logrus.InfoLevel, "Matched incoming DID..")
		valid, err := h.userStore.FinishValidation(number, info.DidWorkspaceId)
		if err != nil {
			return utils.HandleInternalErr("IncomingDIDValidation error 2 valid", err, c)
		}
		if !valid {
			return utils.HandleInternalErr("number not valid..", err, c)
		}
		return c.String(http.StatusOK, "network_managed")
	}

	utils.Log(logrus.InfoLevel, "error in number lookup: "+err.Error())

	utils.Log(logrus.InfoLevel, "Looking up BYO DIDs now...")
	byoInfo, err := h.userStore.IncomingBYODIDValidation(did)
	if err == nil {
		match, err := h.userStore.CheckBYOPSTNIPWhitelist(did, source)
		if err != nil {
			return utils.HandleInternalErr("IncomingDIDValidation error 3", err, c)
		}
		if !match {
			return utils.HandleInternalErr("IncomingDIDValidation no match found 2", err, c)
		}
		utils.Log(logrus.InfoLevel, "Matched incoming DID..")
		valid, err := h.userStore.FinishValidation(number, byoInfo.DidWorkspaceId)
		if err != nil {
			return utils.HandleInternalErr("IncomingDIDValidation error 4 valid", err, c)
		}
		if !valid {
			return utils.HandleInternalErr("number not valid..", err, c)
		}

		return c.String(http.StatusOK, "byo_carrier")
	}
	return utils.HandleInternalErr("IncomingDIDValidation error 5", errors.New("No DID match found..."), c)
}

/*
Input: fromdomain
Todo : Looking up SIP Server and find matched one
Output: If success return SIP Ipaddress else return err
*/
func (h *Handler) IncomingTrunkValidation(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "IncomingTrunkValidation is called...")

	//did := c.QueryParam("did")
	//number := c.QueryParam("number")
	//source := c.QueryParam("source")
	fromdomain := c.QueryParam("fromdomain")
	//destDomain := c.QueryParam("destdomain")

	trunkip, err := utils.LookupSIPAddress(fromdomain)
	if err != nil {
		return utils.HandleInternalErr("IncomingTrunkValidation error 4 valid", err, c)
	}

	utils.Log(logrus.InfoLevel, fmt.Sprintf("From domain %s trunk IP is %s..\r\n", fromdomain, *trunkip))

	result, err := h.userStore.IncomingTrunkValidation(*trunkip)
	if err != nil {
		return utils.HandleInternalErr("IncomingTrunkValidation error 1 valid", err, c)
	}

	if result == nil {
		return utils.HandleInternalErr("Checked all SIP trunks no matches were found... error.", errors.New("no matches"), c)
	}
	return c.JSONBlob(http.StatusOK, result)
}

/*
Input: fromdomain
Todo : Looking up SIP Server and find matched one
Output: If success return SIP Ipaddress else return err
*/
func (h *Handler) LookupSIPTrunkByDID(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "LookupSIPTrunkByDID is called...")

	did := c.QueryParam("did")

	result, err := h.userStore.LookupSIPTrunkByDID(did)
	if err != nil {
		return utils.HandleInternalErr("LookupSIPTrunkByDID error", err, c)
	}

	if result == nil {
		return utils.HandleInternalErr("checked all SIP trunks and found that none were online... error.", errors.New("no SIP"), c)
	}

	return c.JSONBlob(http.StatusOK, result)
}

/*
Input: source
Todo : Looking up MediaServer and find matched one
Output: If success return StatusNoContent else return err
*/
func (h *Handler) IncomingMediaServerValidation(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "IncomingMediaServerValidation is called...")

	//number:= c.QueryParam("number")
	source := c.QueryParam("source")
	//did := c.QueryParam("did")

	result, err := h.userStore.IncomingMediaServerValidation(source)

	if err != nil {
		return utils.HandleInternalErr("IncomingMediaServerValidation error", err, c)
	}

	if result {
		return c.NoContent(http.StatusNoContent)
	}
	return utils.HandleInternalErr("No media server found", errors.New("no media server"), c)
}

/*
Input: domain, user
Todo : Update extensions with domain, user, workspace
Output: If success return StatusNoContent else return err
*/
func (h *Handler) StoreRegistration(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "StoreRegistration is called...")

	domain := c.FormValue("domain")
	//ip := rc.FormValue("ip")
	user := c.FormValue("user")
	//contact := c.FormValue("contact")
	workspace, err := h.callStore.GetWorkspaceByDomain(domain)
	if err != nil {
		return utils.HandleInternalErr("StoreRegistration error..", err, c)
	}

	expires, err := strconv.Atoi(c.FormValue("expires"))
	if err != nil {
		utils.Log(logrus.InfoLevel, "Could not get expiry.. not setting online\r\n")
		return c.NoContent(http.StatusOK)
	}

	err = h.userStore.StoreRegistration(user, expires, workspace)
	if err != nil {
		return utils.HandleInternalErr("StoreRegistration Could not execute query..", err, c)
	}
	return c.NoContent(http.StatusOK)
}

/*
Input: domain, user
Todo : Get settings
Output: If success return Settings model else return err
*/
func (h *Handler) GetSettings(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetSettings")

	settings, err := h.userStore.GetSettings()

	if err != nil {
		return utils.HandleInternalErr("GetSettings error:"+err.Error(), err, c)
	}
	return c.JSON(http.StatusOK, &settings)
}

/*
Input: did
Todo : Get SIP URI with matching did number
Output: If success return sip uri else return err
*/
func (h *Handler) ProcessSIPTrunkCall(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "ProcessSIPTrunkCall is called")

	did := c.QueryParam("did")

	result, err := h.userStore.ProcessSIPTrunkCall(did)
	if err != nil {
		return utils.HandleInternalErr("ProcessSIPTrunkCall error valid", err, c)
	}

	return c.JSONBlob(http.StatusOK, result)
}

/*
Input: did
Todo : Get dialplan route
Output: If success return dialpan context
*/
func (h *Handler) ProcessDialplan(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "ProcessDialplan is called")

	//username := c.QueryParam("username")
	//domain := c.QueryParam("domain")
	//routerip := c.QueryParam("routerip")
	requestUser := c.QueryParam("requestuser")
	//addr := c.QueryParam("addr")

	result, err := h.userStore.ProcessDialplan(requestUser)
	if err != nil {
		return utils.HandleInternalErr("ProcessDialplan error valid", err, c)
	}

	return c.JSONBlob(http.StatusOK, result)
}

/*
Input: sip_call_id
Output: If success return nil
*/
func (h *Handler) ProcessCDRsAndBill(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "ProcessCDRsAndBill is called")

	sipCallId := c.FormValue("callid")


	call, err := h.callStore.GetCallBySIPCallId(sipCallId)
	if err != nil {
		utils.Log(logrus.InfoLevel, fmt.Sprintf("ProcessCDRsAndBill couldnt find call %s in records -- skipping the process routine and sending back 2XX anyway", sipCallId));
		return c.NoContent(http.StatusOK);
		//return utils.HandleInternalErr("ProcessCDRsAndBill error while looking up call.", err, c)
	}


	utils.Log(logrus.InfoLevel, fmt.Sprintf("ProcessCDRsAndBill looking up workspace %d", call.WorkspaceId))
	workspace, err := h.callStore.GetWorkspaceFromDB(call.WorkspaceId)
	if err != nil {
		return utils.HandleInternalErr("ProcessCDRsAndBill error looking up workspace.", err, c)
	}

	callStartDate, err := utils.ParseDate(call.CreatedAt)
	if err != nil {
		return utils.HandleInternalErr("ProcessCDRsAndBill error parsing date string.", err, c)
	}

	seconds := utils.CalculateCallDuration( callStartDate )

	debit := model.Debit{
		Source: "CALL",
		Status: "completed",
		Seconds: seconds,
		ModuleId: call.Id,
		UserId: call.UserId,
	}

	// Get Call Rate depends number and type
	rate := utils.LookupBestCallRate2(call.From, call.To, debit.Type)
	if rate == nil {
		return c.NoContent(http.StatusNotFound)
	}

	// update the calls status to ended
	update := model.CallUpdate{CallId: call.Id, Status: "ended"}
	err = h.callStore.UpdateCall(&update)
	if err != nil {
		return utils.HandleInternalErr("UpdateCall Could not execute query..", err, c)
	}

	debit.PlanSnapshot = workspace.Plan

	err = h.debitStore.CreateDebit(rate, &debit)
	if err != nil {
		return utils.HandleInternalErr("ProcessCDRsAndBill could not create debit.", err, c)
	}

	// send CDR to any remote locations configured by the user
	err = utils.CreateCDRs(call)
	if err != nil {
		return utils.HandleInternalErr("ProcessCDRsAndBill could not create CDRs", err, c)
	}

	return c.NoContent(http.StatusCreated);
}

/*
Input: sip_msg
Todo : Get SIP URI with matching did number
Output: If success return sip uri else return err
*/
func (h *Handler) CaptureSIPMessage(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "CaptureSIPMessage	s called")

	domain := c.QueryParam("domain")
	sipMsg := c.QueryParam("sip_msg")

	result, err := h.userStore.CaptureSIPMessage(domain, sipMsg)
	if err != nil {
		return utils.HandleInternalErr("ProcessSIPTrunkCall error valid", err, c)
	}

	return c.JSONBlob(http.StatusOK, result)
}

/*
Input: invite_ip
Todo : Post processing for invite call event
Output: If success return sip uri else return err
*/
func (h *Handler) LogCallInviteEvent(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "LogCallInviteEvent	s called")

	inviteIp := c.QueryParam("invite_ip")

	err := h.userStore.LogCallInviteEvent(inviteIp)
	if err != nil {
		return utils.HandleInternalErr("LogCallInviteEvent error valid", err, c)
	}

	return c.NoContent(http.StatusOK)
}

/*
Input: invite_ip
Todo : Post processing for invite call event
Output: If success return sip uri else return err
*/
func (h *Handler) LogCallByeEvent(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "LogCallByeEvent	s called")

	inviteIp := c.QueryParam("invite_ip")

	err := h.userStore.LogCallByeEvent(inviteIp)
	if err != nil {
		return utils.HandleInternalErr("LogCallByeEvent error valid", err, c)
	}

	return c.NoContent(http.StatusOK)
}
