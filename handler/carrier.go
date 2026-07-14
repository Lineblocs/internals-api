package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/helpers"
	"lineblocs.com/api/utils"
)

/*
Input: callid, status
Todo : Update sip_status of calls with matching sip_call_id
Output: If success return NoContent else return err
*/
func (h *Handler) CreateSIPReport(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "CreateSIPReport is called...")

	payload := map[string]interface{}{}
	if err := c.Bind(&payload); err != nil && !errors.Is(err, echo.ErrUnsupportedMediaType) {
		return c.JSON(http.StatusBadRequest, "invalid request body")
	}

	getString := func(keys ...string) string {
		for _, key := range keys {
			if value, ok := payload[key]; ok {
				switch v := value.(type) {
				case string:
					if strings.TrimSpace(v) != "" {
						return strings.TrimSpace(v)
					}
				case float64:
					return strconv.FormatInt(int64(v), 10)
				case int:
					return strconv.Itoa(v)
				}
			}

			formValue := strings.TrimSpace(c.FormValue(key))
			if formValue != "" {
				return formValue
			}

			queryValue := strings.TrimSpace(c.QueryParam(key))
			if queryValue != "" {
				return queryValue
			}
		}

		return ""
	}

	callid := getString("callid", "call_id", "sip_call_id")
	status := getString("status", "sip_status")

	if callid == "" || status == "" {
		return c.JSON(http.StatusBadRequest, "callid and status are required")
	}
	if _, err := strconv.Atoi(status); err != nil {
		return c.JSON(http.StatusBadRequest, "status must be numeric")
	}

	err := h.carrierStore.CreateSIPReport(callid, status)
	if err != nil {
		return utils.HandleInternalErr("CreateSIPReport error", err, c)
	}
	return c.NoContent(http.StatusOK)
}

/*
Input: callto, callfrom, userid
Todo : Create and Start Router Flow
Output: If success return host Ipaddress else return err
*/
func (h *Handler) ProcessRouterFlow(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "ProcessRouterFlow is called...")

	var flow *helpers.Flow
	callto := c.QueryParam("callto")
	callfrom := c.QueryParam("callfrom")
	userId := c.QueryParam("userid")

	destCode, err := helpers.ParseCountryCode(callto)

	if err != nil {
		return utils.HandleInternalErr("ProcessRouterFlow invalid destination number", err, c)
	}
	utils.Log(logrus.InfoLevel, fmt.Sprintln("Dest Code is: "+destCode))

	originCode, err := helpers.ParseCountryCode(callfrom)

	if err != nil {
		return utils.HandleInternalErr("ProcessRouterFlow invalid source number", err, c)
	}

	utils.Log(logrus.InfoLevel, fmt.Sprintln("Source Code is: "+destCode))

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
		return utils.HandleInternalErr("ProcessRouterFlow error 2", err, c)
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
