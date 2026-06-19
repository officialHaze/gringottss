package main

import (
	"log"
	"os"
	"path"
	"strings"

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
	args := os.Args[1:]

	if len(args) <= 0 {
		logger.ERROR().Println(`
			Missing Arguments (need either of the 2):
				1. start_server - Starts the Gringottss API server.
				2. build_db - Builds the DB. Recommended to run once before running the program for first time.
		`)
		os.Exit(1)
	}

	arg := args[0]
	switch strings.ToLower(strings.TrimSpace(arg)) {
	case "start_server":
		// start the server
		startServer()
	case "build_db":
		// build the DB
		buildDB()
	}
}

func buildDB() *sqliteclient.Client {
	// Initialize the client and open DB
	db := path.Join("data", os.Getenv("DB_NAME"))
	client, err := sqliteclient.Init("sqlite", db, 1).OpenDB()
	if err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
		return nil
	}

	// Run migrations
	if err := client.RunMigrations(); err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
		return nil
	}

	return client
}

func startServer() {
	client := buildDB()

	// Load sqlc queries for global use
	if err := client.LoadQueries(); err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
		return
	}

	// Start the server
	logger.INFO().Printf("=> Starting GRINGOTTSS API Server on PORT - %d", settings.MySettings.SERVER_PORT)
	server.Start(client)
}
