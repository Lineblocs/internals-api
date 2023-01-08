package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"lineblocs.com/api/helpers"
	"lineblocs.com/api/utils"
)

/*
Input: callid, status
Todo : Update sip_status of calls with matching sip_call_id
Output: If success return NoContent else return err
*/
func (h *Handler) CreateSIPReport(c echo.Context) error {
	callid := c.FormValue("callid")
	status := c.FormValue("status")

	err := h.carrierStore.CreateSIPReport(callid, status)
	if err != nil {
		return utils.HandleInternalErr("CreateSIPReport error", err, c)
	}
	return c.NoContent(http.StatusNoContent)
}

/*
Input: callto, callfrom, userid
Todo : Create and Start Router Flow
Output: If success return host Ipaddress else return err
*/
func (h *Handler) ProcessRouterFlow(c echo.Context) error {
	var flow *helpers.Flow
	callto := c.Param("callto")
	callfrom := c.Param("callfrom")
	userId := c.Param("userid")

	destCode, err := helpers.ParseCountryCode(callto)

	if err != nil {
		panic(err)
	}
	fmt.Println("code is: " + destCode)

	originCode, err := helpers.ParseCountryCode(callfrom)

	if err != nil {
		panic(err)
	}
	fmt.Println("code is: " + originCode)

	// Lookup flow or country flow
	flow, err = h.carrierStore.CreateRoutingFlow(&callfrom, &callto, &userId)
	if err != nil {
		return utils.HandleInternalErr("ProcessRouterFlow error 1", err, c)
	}

	data := make(map[string]string)
	data["origin_code"] = originCode
	data["dest_code"] = destCode
	data["from"] = callfrom
	data["to"] = callto

	// Start processing flow with helpers
	providers, err := h.carrierStore.StartProcessingFlow(flow, data)

	if err != nil {
		panic(err)
	}
	if len(providers) == 0 {
		return utils.HandleInternalErr("No providers available..", err, c)
	}
	provider := providers[0]
	if len(provider.Hosts) == 0 {
		return utils.HandleInternalErr("No IPs to route to..", err, c)
	}
	host := provider.Hosts[0]
	return c.JSON(http.StatusOK, []byte(host.IPAddr))
}
