package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	helpers "github.com/Lineblocs/go-helpers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/mrwaggel/golimiter"
	"github.com/sirupsen/logrus"
	"lineblocs.com/api/database"
	"lineblocs.com/api/handler"
	"lineblocs.com/api/model"
	"lineblocs.com/api/router"
	"lineblocs.com/api/store"
	"lineblocs.com/api/utils"
)

// Global service variables
var dbConn *database.MySQLConn
var rdb *redis.Client
var cqlCluster *gocql.ClusterConfig
var cqlSess *gocql.Session
var data *model.ServerData
var customizations *helpers.CustomizationSettings

// updateCustomizationSettings fetches configuration updates from helper utilities.
// Any retrieval failure is logged and returned as an error to be handled by the caller.
func updateCustomizationSettings() error {
	record, err := helpers.GetCustomizationSettings()
	if err != nil {
		utils.Log(logrus.ErrorLevel, err.Error())
		return err
	}

	customizations = record
	return nil
}

func main() {
	// Initialize Logrus and configure targets based on configurations
	logDestination := utils.Config("LOG_DESTINATIONS")
	helpers.InitLogrus(logDestination)

	utils.Log(logrus.InfoLevel, "Running setup methods for api server..")
	
	// Fetch and load media_server listing. Clean exit on error to allow orchestrator (K8s/Docker) to restart.
	servers, err := helpers.CreateMediaServers()
	if err != nil {
		utils.Log(logrus.ErrorLevel, fmt.Sprintf("Failed to create media servers: %v", err))
		os.Exit(1)
	}

	data = &model.ServerData{
		Mutex:   sync.RWMutex{},
		Servers: servers,
	}

	// Establish connection with the MySQL database.
	utils.Log(logrus.InfoLevel, "Connecting to database...")
	db, err := helpers.CreateDBConn()
	if err != nil {
		utils.Log(logrus.ErrorLevel, fmt.Sprintf("Failed to connect to database: %v", err))
		os.Exit(1)
	}

	dbConn = database.NewMySQLConn(db)

	// Establish connection with the Redis cache.
	rdb, err = helpers.CreateRedisConn()
	if err != nil {
		utils.Log(logrus.ErrorLevel, fmt.Sprintf("Failed to connect to Redis: %v", err))
		os.Exit(1)
	}

	// Retrieve primary configuration records prior to triggering routing.
	err = updateCustomizationSettings()
	if err != nil {
		utils.Log(logrus.ErrorLevel, fmt.Sprintf("Failed to fetch initial customizations: %v", err))
		os.Exit(1)
	}

	utils.Log(logrus.InfoLevel, fmt.Sprintf("got customization settings. billing frequency = %s", customizations.BillingFrequency))

	// Setup context for background routines to stop cleanly during shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize background sync intervals
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// Spin off asynchronous settings-update routine
	go func() {
		for {
			select {
			case <-ticker.C:
				utils.Log(logrus.InfoLevel, "updating customization settings")
				err := updateCustomizationSettings()
				if err != nil {
					utils.Log(logrus.DebugLevel, "failed to get updated customizations record")
					continue
				}
				utils.Log(logrus.InfoLevel, "updated settings successfully.")
			case <-ctx.Done():
				utils.Log(logrus.InfoLevel, "stopping customization settings updater routine")
				return
			}
		}
	}()

	// Construct HTTP router framework and inject data stores
	r := buildServer()

	// Spin up Echo engine in a non-blocking goroutine
	go func() {
		utils.Log(logrus.InfoLevel, "Starting API...")
		if utils.Config("USE_TLS") == "on" {
			certPath := utils.Config("TLS_CERT_PATH")
			keyPath := utils.Config("TLS_KEY_PATH")
			httpsPort := utils.ReadEnv("HTTPS_PORT", "443")
			utils.Log(logrus.InfoLevel, fmt.Sprintf("Starting HTTP server with TLS. cert=%s, key=%s\r\n", certPath, keyPath))
			if err := r.StartTLS(":"+httpsPort, certPath, keyPath); err != nil && err != http.ErrServerClosed {
				utils.Log(logrus.ErrorLevel, fmt.Sprintf("TLS server error: %v", err))
				os.Exit(1)
			}
		} else {
			httpPort := utils.ReadEnv("HTTP_PORT", "80")
			utils.Log(logrus.InfoLevel, fmt.Sprintf("HTTP port %s\r\n", httpPort))
			if err := r.Start(":" + httpPort); err != nil && err != http.ErrServerClosed {
				utils.Log(logrus.ErrorLevel, fmt.Sprintf("HTTP server error: %v", err))
				os.Exit(1)
			}
		}
	}()

	// Monitor OS termination signals to coordinate shutdown
	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGINT, syscall.SIGTERM)
	
	// Wait for termination signal
	<-shutdownSig
	utils.Log(logrus.InfoLevel, "Shutting down service gracefully...")

	// Cancel background routines immediately
	cancel()

	// Implement graceful shutdown with safety timeout context (10 seconds)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := r.Shutdown(shutdownCtx); err != nil {
		utils.Log(logrus.ErrorLevel, fmt.Sprintf("Server forced to shutdown: %v", err))
		os.Exit(1)
	}

	utils.Log(logrus.InfoLevel, "Service stopped cleanly.")
	os.Exit(0)
}

// buildServer constructs the Echo framework, creates stores, and configures route targets
func buildServer() *echo.Echo {
	utils.SetSetting(model.GlobalSettings{ValidateCallerId: false})

	r := router.New()
	utils.Log(logrus.InfoLevel, "Configuring HTTP server...")
	
	if utils.Config("USE_LIMIT_MIDDLEWARE") == "on" {
		r.Any("", limitHandler)
	}

	// Instantiate operational store nodes
	as := store.NewAdminStore(dbConn)
	cs := store.NewCallStore(dbConn)
	crs := store.NewCarrierStore(dbConn)
	ds := store.NewDebitStore(dbConn)
	fs := store.NewFaxStore(dbConn)
	ls := store.NewLoggerStore(dbConn)
	rs := store.NewRecordingStore(dbConn)
	us := store.NewUserStore(dbConn, rdb)
	h := handler.NewHandler(as, cs, crs, ds, fs, ls, rs, us)

	// Bind endpoints
	h.Register(r)
	return r
}

// limitHandler intercepts API calls to prevent brute force or denial attempts
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

	var limit int = 60
	if isCarrier {
		limit = 3600
	}

	var indexLimiter = golimiter.New(limit, time.Minute)

	if indexLimiter.IsLimited(addr) {
		return c.String(http.StatusTooManyRequests, fmt.Sprintf("Rate limit exhausted from %s", addr))
	}
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