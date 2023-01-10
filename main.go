package main

import (
	"net/http"
	"os"
	"sync"
	"time"

	//"errors"
	"database/sql"
	"fmt"

	lineblocs "github.com/Lineblocs/go-helpers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/mrwaggel/golimiter"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/handler"
	"lineblocs.com/api/model"
	"lineblocs.com/api/router"
	"lineblocs.com/api/store"
	"lineblocs.com/api/utils"
)

var db *sql.DB
var data *model.ServerData

func main() {
	utils.InitLogrus()

	utils.Log(logrus.InfoLevel, "Starting API...")

	// Load media_server list from db and create media server
	var err error
	servers, err := lineblocs.CreateMediaServers()

	data = &model.ServerData{
		Mutex:   sync.RWMutex{},
		Servers: servers}

	if err != nil {
		utils.Log(logrus.PanicLevel, err.Error())
		panic(err)
	}

	// Create DB Connection with MySQL
	db, err = lineblocs.CreateDBConn()
	if err != nil {
		utils.Log(logrus.PanicLevel, err.Error())
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		// Start Internals-API Backend server
		startServer()
		wg.Done()
	}()
	wg.Wait()

}

// Start Internals-API Backend server
// Configure Handler, limit middleware, TLS
func startServer() {
	utils.SetSetting(model.GlobalSettings{ValidateCallerId: false})

	// Start Server with Echo
	r := router.New()
	utils.Log(logrus.InfoLevel, "Starting HTTP server...")
	// Configure Limit Handler if USE_LIMIT_MIDDLEWARE is "on"
	if os.Getenv("USE_LIMIT_MIDDLEWARE") == "on" {
		r.Any("", limitHandler)
	}

	// Configure Handler with Global DB
	as := store.NewAdminStore(db)
	cs := store.NewCallStore(db)
	crs := store.NewCarrierStore(db)
	ds := store.NewDebitStore(db)
	fs := store.NewFaxStore(db)
	ls := store.NewLoggerStore(db)
	rs := store.NewRecordingStore(db)
	us := store.NewUserStore(db)
	h := handler.NewHandler(as, cs, crs, ds, fs, ls, rs, us)

	// Register Handler for Echo context
	h.Register(r)

	// Start with 443 port if TLS is ON
	utils.Log(logrus.InfoLevel, "Starting HTTP server without TLS\r\n")
	if os.Getenv("USE_TLS") == "on" {
		certPath := os.Getenv("TLS_CERT_PATH")
		keyPath := os.Getenv("TLS_KEY_PATH")

		utils.Log(logrus.InfoLevel, fmt.Sprintf("Starting HTTP server with TLS. cert=%s,  key=%s\r\n", certPath, keyPath))
		r.Logger.Fatal(r.StartTLS(":443", certPath, keyPath))
		utils.Log(logrus.InfoLevel, "Started server...")
		return
	}

	// Start with 80 port if TLS is OFF
	httpPort := utils.ReadEnv("HTTP_PORT", "80")
	utils.Log(logrus.InfoLevel, fmt.Sprintf("HTTP port %s\r\n", httpPort))
	r.Logger.Fatal(r.Start(":" + httpPort))
	utils.Log(logrus.InfoLevel, "Started server...")
}

// Configure Limit Handler for Echo context
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

	// Limit for users

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
