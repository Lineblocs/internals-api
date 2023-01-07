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
	"lineblocs.com/api/handler"
	"lineblocs.com/api/router"
	"lineblocs.com/api/store"
	"lineblocs.com/api/utils"
)

type GlobalSettings struct {
	ValidateCallerId bool
}

type ServerData struct {
	mu      sync.RWMutex
	servers []*lineblocs.MediaServer
}

var db *sql.DB
var settings *GlobalSettings
var data *ServerData

func main() {
	fmt.Println("starting API...")
	var err error
	servers, err := lineblocs.CreateMediaServers()

	data = &ServerData{
		mu:      sync.RWMutex{},
		servers: servers}

	if err != nil {
		panic(err)
	}
	db, err = lineblocs.CreateDBConn()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		startServer()
		wg.Done()
	}()
	wg.Wait()

}

func startServer() {
	settings = &GlobalSettings{ValidateCallerId: false}
	r := router.New()
	fmt.Println("starting HTTP server...")

	if os.Getenv("USE_LIMIT_MIDDLEWARE") == "on" {
		r.Any("", limitHandler)
	}

	as := store.NewAdminStore(db)
	cs := store.NewCallStore(db)
	crs := store.NewCarrierStore(db)
	ds := store.NewDebitStore(db)
	fs := store.NewFaxStore(db)
	ls := store.NewLoggerStore(db)
	rs := store.NewRecordingStore(db)
	us := store.NewUserStore(db)
	h := handler.NewHandler(as, cs, crs, ds, fs, ls, rs, us)
	h.Register(r)

	fmt.Printf("Starting HTTP server without TLS\r\n")
	if os.Getenv("USE_TLS") == "on" {
		certPath := os.Getenv("TLS_CERT_PATH")
		keyPath := os.Getenv("TLS_KEY_PATH")

		fmt.Printf("Starting HTTP server with TLS. cert=%s,  key=%s\r\n", certPath, keyPath)
		r.Logger.Fatal(r.StartTLS(":443", certPath, keyPath))
		fmt.Println("started server...")
		return
	}

	httpPort := utils.ReadEnv("HTTP_PORT", "80")
	fmt.Printf("HTTP port %s\r\n", httpPort)
	r.Logger.Fatal(r.Start(":" + httpPort))
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
