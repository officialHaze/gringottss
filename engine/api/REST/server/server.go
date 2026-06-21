package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/officialhaze/gringottss/api-server/api/REST/server/routes"
	"github.com/officialhaze/gringottss/api-server/logger"
	"github.com/officialhaze/gringottss/api-server/settings"
	"github.com/officialhaze/gringottss/api-server/sqliteclient"
	"github.com/officialhaze/gringottss/api-server/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start(sqliteClient *sqliteclient.Client) {
	if !util.InDevMode() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.SetTrustedProxies([]string{
		"10.0.0.0/8",
		"172.16.0.0/12", // covers all Docker default subnets
		"192.168.0.0/16",
	})

	// Cors setup
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Type"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// Map routes
	routes.MapRoutes(r)

	addr := fmt.Sprintf(":%d", settings.MySettings.SERVER_PORT)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Run the server in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ERROR().Printf("Listen: %s", err.Error())
			os.Exit(1)
		}
	}()

	// Listen for close signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println()
	logger.INFO().Println("Shutting down server gracefully...")

	// Give in-flight requests some time to complete
	timeout := 10 * time.Second
	logger.INFO().Println("Waiting for any in-flight request(s) to complete...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// // Disconnect from redis
	// if err := redisclient.Client.Close(); err != nil {
	// 	logger.ERROR().Println(err.Error())
	// }
	// logger.SUCCESS().Println("Disconnected from Redis!")

	// Close the DB
	logger.INFO().Println("Closing DB...")
	if err := sqliteClient.DB.Close(); err != nil {
		logger.ERROR().Printf("Failed to close DB!\n%s", err.Error())
	} else {
		logger.SUCCESS().Println("DB closed successfully!")
	}

	// Shut down server
	if err := srv.Shutdown(ctx); err != nil {
		logger.ERROR().Println("Server forced to shutdown!")
		os.Exit(1)
	}

	logger.SUCCESS().Println("Server has been shutdown gracefully!")
}
