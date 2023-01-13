package handler

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"lineblocs.com/api/utils"
)

/*
Register Routers here
Matching API end points and Handler function
*/

func (h *Handler) Register(r *echo.Echo) {
	if os.Getenv("USE_AUTH_MIDDLEWARE") == "on" {
		// Set BasicAuth Middleware
		r.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			if h.userStore.ValidateAccess(username, password) {
				utils.SetMicroservice(username)
				return true, nil
			}
			return false, nil
		}))
	}

	// Call Related Routing
	r.POST("/call/createCall", h.CreateCall)
	r.POST("/call/updateCall", h.UpdateCall)
	r.GET("/call/fetchCall", h.FetchCall)
	r.POST("/call/setSIPCallID", h.SetSIPCallID)
	r.POST("/call/setProviderByIP", h.SetProviderByIP)
	r.POST("/conference/createConference", h.CreateConference)

	// Debit Related Routing
	r.POST("/debit/createDebit", h.CreateDebit)
	r.POST("/debit/createAPIUsageDebit", h.CreateAPIUsageDebit)

	// Debugger Log Related Routing
	r.POST("/debugger/createLog", h.CreateLog)
	r.POST("/debugger/createLogSimple", h.CreateLogSimple)

	// Fax Related Routing
	r.POST("/fax/createFax", h.CreateFax)

	// Recording Related Routing
	r.POST("/recording/createRecording", h.CreateRecording)
	r.POST("/recording/updateRecording", h.UpdateRecording)
	r.POST("/recording/updateRecordingTranscription", h.UpdateRecordingTranscription)
	r.GET("/recording/getRecording", h.GetRecording)

	// Carrier Related Routing
	r.POST("/carrier/createSIPReport", h.CreateSIPReport)
	r.GET("/carrier/processRouterFlow", h.ProcessRouterFlow)

	// User Related Routing
	r.GET("/user/verifyCaller", h.VerifyCaller)
	r.GET("/user/verifyCallerByDomain", h.VerifyCallerByDomain)
	r.GET("/user/getUserByDomain", h.GetUserByDomain)
	r.GET("/user/getUserByDID", h.GetUserByDID)
	r.GET("/user/getUserByTrunkSourceIp", h.GetUserByTrunkSourceIp)
	r.GET("/user/getWorkspaceMacros", h.GetWorkspaceMacros)
	r.GET("/user/getDIDNumberData", h.GetDIDNumberData)
	r.GET("/user/getPSTNProviderIP", h.GetPSTNProviderIP)
	r.GET("/user/getPSTNProviderIPForTrunk", h.GetPSTNProviderIPForTrunk)
	r.GET("/user/ipWhitelistLookup", h.IPWhitelistLookup)
	r.GET("/user/getDIDAcceptOption", h.GetDIDAcceptOption)
	r.GET("/user/getDIDAssignedIP", h.GetDIDAssignedIP)
	r.GET("/user/getUserAssignedIP", h.GetUserAssignedIP)
	r.GET("/user/getTrunkAssignedIP", h.GetTrunkAssignedIP)
	r.GET("/user/addPSTNProviderTechPrefix", h.AddPSTNProviderTechPrefix)
	r.GET("/user/getCallerIdToUse", h.GetCallerIdToUse)
	r.GET("/user/getExtensionFlowInfo", h.GetExtensionFlowInfo)
	r.GET("/user/getFlowInfo", h.GetFlowInfo)
	r.GET("/user/getDIDDomain", h.GetDIDDomain)
	r.GET("/user/getCodeFlowInfo", h.GetCodeFlowInfo)
	r.GET("/user/incomingDIDValidation", h.IncomingDIDValidation)
	r.GET("/user/incomingTrunkValidation", h.IncomingTrunkValidation)
	r.GET("/user/lookupSIPTrunkByDID", h.LookupSIPTrunkByDID)
	r.GET("/user/incomingMediaServerValidation", h.IncomingMediaServerValidation)
	r.POST("/user/storeRegistration", h.StoreRegistration)
	r.GET("/user/getSettings", h.GetSettings)
	r.GET("/user/processSIPTrunkCall", h.ProcessSIPTrunkCall)

	// Admin Related Routing
	r.GET("/healthz", h.Healthz)
	r.POST("/admin/sendAdminEmail", h.SendAdminEmail)
	r.GET("/getBestRTPProxy", h.GetBestRTPProxy)

}
