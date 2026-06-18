package main

import (
	"log"
	"os"

	"github.com/officialhaze/gringottss/api-server/api/REST/server"
	"github.com/officialhaze/gringottss/api-server/logger"
	"github.com/officialhaze/gringottss/api-server/settings"
	"github.com/officialhaze/gringottss/api-server/sqliteclient"
	"github.com/officialhaze/gringottss/api-server/util"

	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	// Generate my settings
	if err := settings.Generate(); err != nil {
		log.Fatalln(err)
	}

	// load the env file
	if err := godotenv.Load(settings.MySettings.ENV_FILE_NAME); err != nil {
		log.Fatalln(err)
	}

	// Setup lumberjack and logger
	logfile := &lumberjack.Logger{
		Filename:   "/var/log/gringottss-server.log",
		MaxSize:    100, // mb
		MaxBackups: 5,
		MaxAge:     28, // days
		Compress:   true,
	}

	// Log on both console and logfile when in dev mode
	if util.InDevMode() {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(logfile)
	}
}

func main() {
	// Initialize the client and open DB
	client, err := sqliteclient.Init("sqlite", "gringottss.db", 1).OpenDB()
	if err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
	}

	// Run migrations
	if err := client.RunMigrations(); err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
	}

	// Load sqlc queries for global use
	if err := client.LoadQueries(); err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
	}

	// Start the server
	logger.INFO().Printf("=> Starting GRINGOTTSS API Server on PORT - %d", settings.MySettings.SERVER_PORT)
	server.Start(client)
}
