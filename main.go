package main

import (
	"net/http"
	"sync"
	"time"
	//"errors"
	"github.com/gocql/gocql"
	"fmt"

	helpers "github.com/Lineblocs/go-helpers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/mrwaggel/golimiter"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/handler"
	"lineblocs.com/api/model"
	"lineblocs.com/api/router"
	"lineblocs.com/api/store"
	"lineblocs.com/api/utils"
	"lineblocs.com/api/database"
)

var dbConn *database.MySQLConn
var rdb *redis.Client
var cqlCluster *gocql.ClusterConfig
var cqlSess *gocql.Session
var data *model.ServerData
var customizations *helpers.CustomizationSettings

func updateCustomizationSettings() (error) {
	record, err := helpers.GetCustomizationSettings()
	if err != nil {
		utils.Log(logrus.PanicLevel, err.Error())
		panic(err)
		return err
	}

	customizations = record
	return nil
}

func main() {
    ticker := time.NewTicker(60 * time.Second)
    defer ticker.Stop()

	// Init Logrus and configure channels
	logDestination := utils.Config("LOG_DESTINATIONS")
	helpers.InitLogrus(logDestination)

	utils.Log(logrus.InfoLevel, "Running setup methods for api server..")
	// Load media_server list from db and create media server
	var err error
	servers, err := helpers.CreateMediaServers()

	data = &model.ServerData{
		Mutex:   sync.RWMutex{},
		Servers: servers}

	if err != nil {
		utils.Log(logrus.PanicLevel, err.Error())
		panic(err)
	}

	// Create DB Connection with MySQL
	utils.Log(logrus.InfoLevel, "Connecting to database...")
	db, err := helpers.CreateDBConn()
	if err != nil {
		utils.Log(logrus.PanicLevel, err.Error())
		panic(err)
	}

	dbConn = database.NewMySQLConn(db)

	rdb, err = helpers.CreateRedisConn()

	// get customization record before starting server
	err = updateCustomizationSettings()
	if err != nil {
		utils.Log(logrus.PanicLevel, err.Error())
		panic(err)
	}

	// connect to cassandra
	/*
	utils.Log(logrus.InfoLevel, "Connecting to cassandra...")
	cassandraAddr := utils.Config("CASSANDRA_HOST") + ":9042"
	cqlCluster = gocql.NewCluster(cassandraAddr)
	cqlCluster.Keyspace = utils.Config("CASSANDRA_KEYSPACE")
	cqlCluster.ProtoVersion = 4
	cqlSess, err = cqlCluster.CreateSession()
	if err != nil {
		utils.Log(logrus.PanicLevel, err.Error())
		panic(err)
	}
	*/

	utils.Log(logrus.InfoLevel, fmt.Sprintf("got customization settings. billing frequency = %s", customizations.BillingFrequency))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
        for {
			utils.Log(logrus.InfoLevel, fmt.Sprintf("updating customization settings"))
            <-ticker.C
            err := updateCustomizationSettings()
			if err != nil {
				utils.Log(logrus.DebugLevel, fmt.Sprintf("failed to get updated customizations record"))
				continue
			}

			utils.Log(logrus.InfoLevel, fmt.Sprintf("updated settings successfully."))
        }
    }()

	go func() {
		// Start Internals-API Backend server
		utils.Log(logrus.InfoLevel, "Starting API...")
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
	if utils.Config("USE_LIMIT_MIDDLEWARE") == "on" {
		r.Any("", limitHandler)
	}

	// Configure Handler with Global DB
	as := store.NewAdminStore(dbConn)
	cs := store.NewCallStore(dbConn)
	crs := store.NewCarrierStore(dbConn)
	ds := store.NewDebitStore(dbConn)
	fs := store.NewFaxStore(dbConn)
	ls := store.NewLoggerStore(dbConn)
	rs := store.NewRecordingStore(dbConn)
	us := store.NewUserStore(dbConn, rdb)
	h := handler.NewHandler(as, cs, crs, ds, fs, ls, rs, us)

	// Register Handler for Echo context
	h.Register(r)

	// Start with 443 port if TLS is ON
	utils.Log(logrus.InfoLevel, "Starting HTTP server without TLS\r\n")
	if utils.Config("USE_TLS") == "on" {
		certPath := utils.Config("TLS_CERT_PATH")
		keyPath := utils.Config("TLS_KEY_PATH")
		httpsPort := utils.ReadEnv("HTTPS_PORT", "443")
		utils.Log(logrus.InfoLevel, fmt.Sprintf("Starting HTTP server with TLS. cert=%s,  key=%s\r\n", certPath, keyPath))
		r.Logger.Fatal(r.StartTLS(":"+httpsPort, certPath, keyPath))
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
	requestedAddr := c.QueryParam("addr")
	if requestedAddr == "" {
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
