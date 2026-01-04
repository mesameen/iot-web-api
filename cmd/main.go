package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mesameen/iot-web-api/src/config"
	"github.com/mesameen/iot-web-api/src/lib/controller"
	"github.com/mesameen/iot-web-api/src/lib/database"
	"github.com/mesameen/iot-web-api/src/lib/service/commands"
	"github.com/mesameen/iot-web-api/src/lib/service/connections"
	"github.com/mesameen/iot-web-api/src/lib/service/device"
	"github.com/mesameen/iot-web-api/src/lib/service/telematics"
	"github.com/mesameen/iot-web-api/src/telemetryservice"
)

func main() {
	// initialize the config
	config.InitConfig()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// initialize telemetry (tracer + logger)
	telem, err := telemetryservice.NewTelemetry(ctx, config.Config.Common.AppName, "v1")
	if err != nil {
		log.Panic(err)
	}
	defer telem.Shutdown(ctx)
	telem.Infof(ctx, "telemetry initialized")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	router.Use(telem.LogRequest())
	router.Use(telem.MeterRequestDuration())
	router.Use(telem.MeterRequestsInFlight())
	apiRouter := router.Group(config.Config.Common.APIRouteGroup)

	db, err := database.New(ctx, telem)
	if err != nil {
		telem.Fatalf(ctx, "Failed to connect to database. Error: %v", err)
	}
	defer db.Close(ctx)
	ctrl := controller.New(telem, db)
	telematicsHandler := telematics.NewHandler(telem, ctrl)
	telematics.RegisterRoutes(apiRouter.Group("/telematics"), telematicsHandler)
	deviceHandler := device.NewHandler(telem, ctrl)
	device.RegisterRoutes(apiRouter.Group("/device"), deviceHandler)
	connectionsHandler := connections.NewHandler(telem, ctrl)
	connections.RegisterRoutes(apiRouter.Group("/connections"), connectionsHandler)
	commandsHandler := commands.NewHandler(telem, ctrl)
	commands.RegisterRoutes(apiRouter.Group("/commands"), commandsHandler)

	server := http.Server{
		Addr:    config.Config.Server.Address,
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			telem.Fatalf(ctx, "Failed to start the server. Error: %v", err)
		}
	}()
	telem.Infof(ctx, "Application is up and running on: %s", config.Config.Server.Address)
	<-ctx.Done()
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 5*time.Second)
	defer timeoutCancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		telem.Fatalf(ctx, "Failed to shutdown the server. Error: %v", err)
		return
	}
	telem.Infof(ctx, "Server shutdown succefully")
}
