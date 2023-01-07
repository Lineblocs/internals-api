package main

import (
	"sync"

	//"errors"
	"database/sql"
	"fmt"

	lineblocs "github.com/Lineblocs/go-helpers"
	_ "github.com/go-sql-driver/mysql"
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

	httpPort := utils.ReadEnv("HTTP_PORT", "80")
	fmt.Printf("HTTP port %s\r\n", httpPort)
	r.Logger.Fatal(r.Start(":" + httpPort))
}
