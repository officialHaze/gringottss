package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"path"
	"strings"

	"github.com/officialhaze/gringottss/api-server/api/REST/server"
	db "github.com/officialhaze/gringottss/api-server/db/sqlc"
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
				3. migrate - Migrates all necessary tables from an existing DB to the current one.
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
	case "migrate":
		if len(args) < 2 {
			logger.ERROR().Println("Not enough arguments. Old DB name must be provided.")
			os.Exit(1)
			return
		}
		dbname := args[1]
		// migrate
		migrate(dbname)
	default:
		logger.ERROR().Println("Invalid/Unsupported argument provided.")
		os.Exit(1)
	}
}

// Build the DB
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

// Start the server
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

// Migrate existing DB tables to the current one
func migrate(dbname string) {
	// Build the current DB (just in case)
	logger.INFO().Println("Building current DB to be sure and avoid crashes during data migration...")
	buildDB()

	// Initialize the client and open old DB
	logger.INFO().Println("Initializing old DB...")
	oldDB := path.Join("migrate", dbname)
	oldDBclient, err := sqliteclient.Init("sqlite", oldDB, 1).OpenDB()
	if err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
		return
	}
	logger.SUCCESS().Println("Old DB initialized.")

	// Run any new current migrations on the old DB just to be safe
	logger.INFO().Println("Running new migrations on old DB...")
	if err := oldDBclient.RunMigrations(); err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
	}
	logger.SUCCESS().Println("Migrations ran successfully on old DB.")

	// Load the queries for old DB
	logger.INFO().Println("Loading queries for old DB...")
	if err := oldDBclient.LoadQueries(); err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
	}
	logger.SUCCESS().Println("Queries loaded for old DB.")

	logger.INFO().Println(" ====> Starting Data Migration <====")
	// Get all the urls in old DB
	urls, err := sqliteclient.Queries.ListURLs(context.Background())
	if err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
	}

	// Iterate over each URL, check if it exists in current DB,
	// Exists? skip, else copy the credentials for the url into current DB.
	for _, url := range urls {
		logger.INFO().Printf("[%s] Checking...", url.Url)

		urlexists := true
		_, err := sqliteclient.Queries.GetURL(context.Background(), url.Url)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				urlexists = false
				break
			}

			// Random error
			logger.ERROR().Printf("[%s] Check failed!\n%s\nSkipping...", url.Url, err.Error())
			continue
		}

		if urlexists {
			logger.WARN().Printf("[%s] Already exists in current DB. Skipping...", url.Url)
			continue
		}

		logger.INFO().Printf("[%s] URL does not exist in current DB. Continuing....", url.Url)

		// Get credentials of the url
		credentials, err := sqliteclient.Queries.GetCredentials(context.Background(), db.GetCredentialsParams{
			Url: url.Url,
		})
		if err != nil {
			logger.ERROR().Printf("[%s] Failed to fetch credentials.\n%s\nSkipping...", url.Url, err.Error())
			continue
		}

		// Copy each credential into current DB
		for i, cred := range credentials {
			logger.INFO().Printf("[%s] Adding credential %d", url.Url, i+1)

			if _, err := sqliteclient.Queries.AddCredential(context.Background(), db.AddCredentialParams{
				ID:             cred.ID,
				Url:            url.Url,
				FormInputID:    cred.FormInputID,
				FormInputName:  cred.FormInputName,
				FormInputXpath: cred.FormInputXpath,
				FormInputType:  cred.FormInputType,
				FormInputVal:   cred.FormInputVal,
			}); err != nil {
				logger.ERROR().Printf("[%s] Failed to add credential %d. Moving to next credential...", url.Url, i+1)
				continue
			}

			logger.SUCCESS().Printf("[%s] Successfully added credential %d", url.Url, i+1)
		}

		logger.SUCCESS().Printf("[%s] Successfully copied to current DB.\n", url.Url)
	}

	logger.SUCCESS().Println("Data migrated successfully!")
}
