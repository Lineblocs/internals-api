package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/utils"
)

/*
Register Routers here
Matching API end points and Handler function
*/

func (h *Handler) Register(e *echo.Echo) {
	g := e.Group("")

	utils.Log(logrus.InfoLevel, "Auth middleware value = "+utils.Config("USE_AUTH_MIDDLEWARE"))
	if utils.Config("USE_AUTH_MIDDLEWARE") == "on" {
		// Set BasicAuth Middleware
		utils.Log(logrus.InfoLevel, "Auth middleware is enabled -- adding API validation")
		g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			if h.userStore.ValidateAccess(username, password) {
				utils.Log(logrus.InfoLevel, "Authentification is successfully passed")
				utils.SetMicroservice(username)
				return true, nil
			}
			return false, nil
		}))
	}

	// For Health Check
	e.GET("/healthz", h.Healthz)

	// Call Related Routing
	g.POST("/call/createCall", h.CreateCall)
	g.POST("/call/updateCall", h.UpdateCall)
	g.GET("/call/fetchCall", h.FetchCall)
	g.POST("/call/setSIPCallID", h.SetSIPCallID)
	g.POST("/call/setProviderByIP", h.SetProviderByIP)
	g.POST("/conference/createConference", h.CreateConference)

	// Debit Related Routing
	g.POST("/debit/createDebit", h.CreateDebit)
	g.POST("/debit/createAPIUsageDebit", h.CreateAPIUsageDebit)

	// Debugger Log Related Routing
	g.POST("/debugger/createLog", h.CreateLog)
	g.POST("/debugger/createLogSimple", h.CreateLogSimple)

	// Fax Related Routing
	g.POST("/fax/createFax", h.CreateFax)

	// Recording Related Routing
	g.POST("/recording/createRecording", h.CreateRecording)
	g.POST("/recording/updateRecording", h.UpdateRecording)
	g.POST("/recording/updateRecordingTranscription", h.UpdateRecordingTranscription)
	g.GET("/recording/getRecording", h.GetRecording)

	// Carrier Related Routing
	g.POST("/carrier/createSIPReport", h.CreateSIPReport)
	g.GET("/carrier/processRouterFlow", h.ProcessRouterFlow)

	// User Related Routing
	g.GET("/user/verifyCaller", h.VerifyCaller)
	g.GET("/user/verifyCallerByDomain", h.VerifyCallerByDomain)
	g.GET("/user/getUserByDomain", h.GetUserByDomain)
	g.GET("/user/getUserByDID", h.GetUserByDID)
	g.GET("/user/getUserByTrunkSourceIp", h.GetUserByTrunkSourceIp)
	g.GET("/user/getWorkspaceMacros", h.GetWorkspaceMacros)
	g.GET("/user/getDIDNumberData", h.GetDIDNumberData)
	g.GET("/user/getPSTNProviderIP", h.GetPSTNProviderIP)
	g.GET("/user/getPSTNProviderIPForTrunk", h.GetPSTNProviderIPForTrunk)
	g.GET("/user/ipWhitelistLookup", h.IPWhitelistLookup)
	g.GET("/user/getDIDAcceptOption", h.GetDIDAcceptOption)
	g.GET("/user/getDIDAssignedIP", h.GetDIDAssignedIP)
	g.GET("/user/getUserAssignedIP", h.GetUserAssignedIP)
	g.GET("/user/getTrunkAssignedIP", h.GetTrunkAssignedIP)
	g.GET("/user/addPSTNProviderTechPrefix", h.AddPSTNProviderTechPrefix)
	g.GET("/user/getCallerIdToUse", h.GetCallerIdToUse)
	g.GET("/user/getExtensionFlowInfo", h.GetExtensionFlowInfo)
	g.GET("/user/getFlowInfo", h.GetFlowInfo)
	g.GET("/user/getDIDDomain", h.GetDIDDomain)
	g.GET("/user/getCodeFlowInfo", h.GetCodeFlowInfo)
	g.GET("/user/incomingDIDValidation", h.IncomingDIDValidation)
	g.GET("/user/incomingTrunkValidation", h.IncomingTrunkValidation)
	g.GET("/user/lookupSIPTrunkByDID", h.LookupSIPTrunkByDID)
	g.GET("/user/incomingMediaServerValidation", h.IncomingMediaServerValidation)
	g.POST("/user/storeRegistration", h.StoreRegistration)
	g.GET("/user/getSettings", h.GetSettings)
	g.GET("/user/processSIPTrunkCall", h.ProcessSIPTrunkCall)
	g.GET("/user/processDialplan", h.ProcessDialplan)
	g.GET("/user/captureSIPMessage", h.CaptureSIPMessage)


	// Admin Related Routing
	g.POST("/admin/sendAdminEmail", h.SendAdminEmail)
	g.GET("/getBestRTPProxy", h.GetBestRTPProxy)

}
