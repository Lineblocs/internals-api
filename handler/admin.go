package handler

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/model"
	"lineblocs.com/api/utils"
)

/*
Todo: Check SIP Routers is Healthy
Output: if success return "OK" else return "Healthz error"
*/
func (h *Handler) Healthz(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "Healthz is called...\r\n")
	err := h.adminStore.Healthz()
	if err != nil {
		return utils.HandleInternalErr("Healthz error", err, c)
	}
	return c.String(http.StatusOK, "OK\n")
}

/*
Input: EmailInfo model
Todo : Send Email to Lineblocs Contact
Output: If success return NoContent else return err
*/
func (h *Handler) SendAdminEmail(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "SendAdminEmail is called...\r\n")
	var emailInfo model.EmailInfo

	if err := c.Bind(&emailInfo); err != nil {
		return utils.HandleInternalErr("SendAdminEmail Could not decode JSON", err, c)
	}
	if err := c.Validate(&emailInfo); err != nil {
		return utils.HandleInternalErr("SendAdminEmail Could not decode JSON", err, c)
	}

	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
	m := mg.NewMessage(
		"Lineblocs <monitor@lineblocs.com>",
		"Admin Error",
		"Admin Error",
		"contact@lineblocs.com")
	body := `<html>
		<head></head>
		<body>
			<h1>Lineblocs Admin Monitor</h1>
			<p>` + emailInfo.Message + `</p>
		</body>
		</html>`

	m.SetHtml(body)
	//m.AddAttachment("files/test.jpg")
	//m.AddAttachment("files/test.txt")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := mg.Send(ctx, m)
	if err != nil {
		return utils.HandleInternalErr("SendAdminEmail error", err, c)
	}
	return c.NoContent(http.StatusNoContent)
}

/*
Input:
Todo : Choose Best one from rtpproxy_sockets
Output: If success return rtpproxy_sock else return "GetBestRTPProxy error"
*/
func (h *Handler) GetBestRTPProxy(c echo.Context) error {
	utils.Log(logrus.InfoLevel, "GetBestRTPProxy is called...\r\n")
	result, err := h.adminStore.GetBestRTPProxy()
	if err != nil {
		return utils.HandleInternalErr("GetBestRTPProxy error", err, c)
	}
	return c.JSON(http.StatusOK, &result)
}
