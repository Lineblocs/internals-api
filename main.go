package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/mrwaggel/golimiter"
	"github.com/sirupsen/logrus"

	helpers "github.com/Lineblocs/go-helpers"
	"lineblocs.com/api/database"
	"lineblocs.com/api/handler"
	"lineblocs.com/api/model"
	"lineblocs.com/api/router"
	"lineblocs.com/api/store"
	"lineblocs.com/api/utils"
)

var dbConn *database.MySQLConn
var rdb *redis.Client
var cqlCluster *gocql.ClusterConfig
var cqlSess *gocql.Session
var data *model.ServerData
var customizations *helpers.CustomizationSettings

func updateCustomizationSettings() error {
	record, err := helpers.GetCustomizationSettings()
	if err != nil {
		logrus.WithError(err).Error("failed to fetch customization settings")
		return err
	}

	customizations = record
	return nil
}

func main() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	shutdownCh := make(chan struct{})

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
		logrus.WithError(err).Error("failed to create media servers")
		return
	}

	// Create DB Connection with MySQL
	utils.Log(logrus.InfoLevel, "Connecting to database...")
	db, err := helpers.CreateDBConn()
	if err != nil {
		logrus.WithError(err).Error("failed to connect to database")
		return
	}

	dbConn = database.NewMySQLConn(db)

	rdb, err = helpers.CreateRedisConn()
	if err != nil {
		logrus.WithError(err).Error("failed to connect to redis")
		return
	}

	// get customization record before starting server
	err = updateCustomizationSettings()
	if err != nil {
		logrus.WithError(err).Error("failed during customization bootstrap")
		return
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

	go func() {
		for {
			select {
			case <-shutdownCh:
				utils.Log(logrus.InfoLevel, "stopping customization refresh loop")
				return
			case <-ticker.C:
			}

			utils.Log(logrus.InfoLevel, fmt.Sprintf("updating customization settings"))
			err := updateCustomizationSettings()
			if err != nil {
				utils.Log(logrus.DebugLevel, fmt.Sprintf("failed to get updated customizations record"))
				continue
			}

			utils.Log(logrus.InfoLevel, fmt.Sprintf("updated settings successfully."))
		}
	}()

	if err := startServer(shutdownCh); err != nil {
		logrus.WithError(err).Error("server exited with error")
	}

}

// Start Internals-API Backend server
// Configure Handler, limit middleware, TLS
func startServer(shutdownCh chan struct{}) error {
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

	httpPort := utils.ReadEnv("HTTP_PORT", "80")
	if utils.Config("USE_TLS") == "on" {
		httpPort = utils.ReadEnv("HTTPS_PORT", "443")
	}

	server := &http.Server{
		Addr:              ":" + httpPort,
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		if utils.Config("USE_TLS") == "on" {
			certPath := utils.Config("TLS_CERT_PATH")
			keyPath := utils.Config("TLS_KEY_PATH")
			utils.Log(logrus.InfoLevel, fmt.Sprintf("starting HTTPS server with TLS. cert=%s key=%s", certPath, keyPath))
			serverErrors <- server.ListenAndServeTLS(certPath, keyPath)
			return
		}

		utils.Log(logrus.InfoLevel, fmt.Sprintf("starting HTTP server on port %s", httpPort))
		serverErrors <- r.StartServer(server)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case sig := <-sigCh:
		close(shutdownCh)
		utils.Log(logrus.InfoLevel, fmt.Sprintf("received shutdown signal: %s", sig.String()))
		shutdownTimeout := utils.ReadEnv("GRACEFUL_SHUTDOWN_TIMEOUT", "10s")
		duration, parseErr := time.ParseDuration(shutdownTimeout)
		if parseErr != nil {
			duration = 10 * time.Second
		}

		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()
		if err := r.Shutdown(ctx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		return nil
	case err := <-serverErrors:
		close(shutdownCh)
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
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
