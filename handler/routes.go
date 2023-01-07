package handler

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mrwaggel/golimiter"
	"lineblocs.com/api/utils"
)

/*
Register Routers here
Use Basic Auth Middleware with Group
*/

func (h *Handler) Register(r *echo.Echo) {

	r.Any("", limitHandler)

	group := r.Group("/", middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	//Call Related Routing
	group.POST("/call/createCall", h.CreateCall)
	group.POST("/call/updateCall", h.UpdateCall)
	group.GET("/call/fetchCall", h.FetchCall)
	group.POST("/call/setSIPCallID", h.SetSIPCallID)
	group.POST("/call/setProviderByIP", h.SetProviderByIP)
	group.POST("/conference/createConference", h.CreateConference)

	//Debit Related Routing
	group.POST("/debit/createDebit", h.CreateDebit)
	group.POST("/debit/createAPIUsageDebit", h.CreateAPIUsageDebit)

	//Debugger Log Related Routing
	group.POST("/debugger/createLog", h.CreateLog)

	//Fax Related Routing
	group.POST("/fax/createFax", h.CreateFax)

	//Recording Related Routing
	group.POST("/recording/createRecording", h.CreateRecording)
	group.POST("/recording/updateRecording", h.UpdateRecording)
	group.POST("/recording/updateRecordingTranscription", h.UpdateRecordingTranscription)
	group.GET("/recording/getRecording", h.GetRecording)

	//Carrier Related Routing
	group.POST("/carrier/createSIPReport", h.CreateSIPReport)
	group.POST("/carrier/processRouterFlow", h.ProcessRouterFlow)

	//User Related Routing
	group.GET("/user/verifyCaller", h.VerifyCaller)
	group.GET("/user/verifyCallerByDomain", h.VerifyCallerByDomain)
	group.GET("/user/getUserByDomain", h.GetUserByDomain)
	group.GET("/user/getUserByDID", h.GetUserByDID)
	group.GET("/user/getUserByTrunkSourceIp", h.GetUserByTrunkSourceIp)
	group.GET("/user/getWorkspaceMacros", h.GetWorkspaceMacros)
	group.GET("/user/getDIDNumberData", h.GetDIDNumberData)
	group.GET("/user/getPSTNProviderIP", h.GetPSTNProviderIP)
	group.GET("/user/getPSTNProviderIPForTrunk", h.GetPSTNProviderIPForTrunk)
	group.GET("/user/ipWhitelistLookup", h.IPWhitelistLookup)
	group.GET("/user/getDIDAcceptOption", h.GetDIDAcceptOption)
	group.GET("/user/getDIDAssignedIP", h.GetDIDAssignedIP)
	group.GET("/user/getUserAssignedIP", h.GetUserAssignedIP)
	group.GET("/user/getTrunkAssignedIP", h.GetTrunkAssignedIP)
	group.GET("/user/addPSTNProviderTechPrefix", h.AddPSTNProviderTechPrefix)
	group.GET("/user/getCallerIdToUse", h.GetCallerIdToUse)
	group.GET("/user/getExtensionFlowInfo", h.GetExtensionFlowInfo)
	group.GET("/user/getFlowInfo", h.GetFlowInfo)
	group.GET("/user/getDIDDomain", h.GetDIDDomain)
	group.GET("/user/getCodeFlowInfo", h.GetCodeFlowInfo)
	group.GET("/user/incomingDIDValidation", h.IncomingDIDValidation)
	group.GET("/user/incomingTrunkValidation", h.IncomingTrunkValidation)
	group.GET("/user/lookupSIPTrunkByDID", h.LookupSIPTrunkByDID)
	group.GET("/user/incomingMediaServerValidation", h.IncomingMediaServerValidation)
	group.GET("/user/storeRegistration", h.StoreRegistration)
	group.GET("/user/getSettings", h.GetSettings)
	group.GET("/user/processSIPTrunkCall", h.ProcessSIPTrunkCall)

	// Admin Related Routing
	group.POST("/admin/sendAdminEmail", h.SendAdminEmail)
	group.GET("/getBestRTPProxy", h.GetBestRTPProxy)

}

func limitHandler(c echo.Context) error {
	var addr string
	requestedAddr := c.Param("addr")
	if &requestedAddr != nil {
		addr = requestedAddr
	} else {
		addr = c.RealIP()
	}

	carrier := c.Request().Header.Get("X-Lineblocs-Carrier-Auth")
	isCarrier := false

	if carrier != "" {
		isCarrier = utils.CheckIfCarrier(carrier)
	}

	// limit for users

	var limit int = 60
	if isCarrier {
		limit = 3600
	}

	var indexLimiter = golimiter.New(limit, time.Minute)

	// Check if the given IP is rate limited
	if indexLimiter.IsLimited(addr) {
		return c.String(http.StatusTooManyRequests, fmt.Sprintf("Rate limit exhausted from %s", addr))
	}
	// Add a request to the count for the Ip
	indexLimiter.Increment(addr)
	totalRequestPastMinute := indexLimiter.Count(addr)
	totalRemaining := limit - totalRequestPastMinute
	return c.String(http.StatusOK, fmt.Sprintf(""+
		"Your IP %s is not rate limited!\n"+
		"You made %d requests in the last minute.\n"+
		"You are allowed to make %d more request.\n"+
		"Maximum request you can make per minute is %d.",
		addr, totalRequestPastMinute, totalRemaining, limit))
}
