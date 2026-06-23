package helpers

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path"

	db "github.com/officialhaze/gringottss/api-server/db/sqlc"
	"github.com/officialhaze/gringottss/api-server/logger"
	"github.com/officialhaze/gringottss/api-server/sqliteclient"
)

// Build the Gringottss DB
func BuildDB() (*sqliteclient.Client, error) {
	// Initialize the client and open DB
	db := path.Join("data", os.Getenv("DB_NAME"))
	client, err := sqliteclient.Init("sqlite", db, 1).OpenDB()
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := client.RunMigrations(); err != nil {
		return nil, err
	}

	return client, nil
}

// Migrate old DB tables to current DB
func MigrateOldDB(dbname string) error {
	// Build the current DB (just in case)
	logger.INFO().Println("Building current DB to be sure and avoid crashes during data migration...")
	newDBClient, err := BuildDB()
	if err != nil {
		logger.ERROR().Println(err.Error())
		return err
	}

	// Load queries for the new DB
	if err := newDBClient.LoadClientQueries(); err != nil {
		logger.ERROR().Println(err.Error())
		return err
	}

	// Initialize the client and open old DB
	logger.INFO().Println("Initializing old DB...")
	oldDB := path.Join("migrate", dbname)
	oldDBclient, err := sqliteclient.Init("sqlite", oldDB, 1).OpenDB()
	if err != nil {
		return err
	}
	logger.SUCCESS().Println("Old DB initialized.")

	// Run any new current migrations on the old DB just to be safe
	logger.INFO().Println("Running new migrations on old DB...")
	if err := oldDBclient.RunMigrations(); err != nil {
		return err
	}
	logger.SUCCESS().Println("Migrations ran successfully on old DB.")

	// Load the queries for old DB
	logger.INFO().Println("Loading queries for old DB...")
	if err := oldDBclient.LoadClientQueries(); err != nil {
		return err
	}
	logger.SUCCESS().Println("Queries loaded for old DB.")

	logger.INFO().Println(" ====> Starting Data Migration <====")
	// Get all the urls in old DB
	urls, err := oldDBclient.Queries.ListURLs(context.Background())
	if err != nil {
		return err
	}

	// Iterate over each URL, check if it exists in current DB,
	// Exists? skip, else copy the credentials for the url into current DB.
	for _, url := range urls {
		logger.INFO().Printf("[%s] Checking...", url.Url)

		urlexists := true
		_, err := newDBClient.Queries.GetURL(context.Background(), url.Url)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				urlexists = false
			} else {
				// Random error
				logger.ERROR().Printf("[%s] Check failed!\n%s\nSkipping...", url.Url, err.Error())
				continue
			}
		}

		if urlexists {
			logger.WARN().Printf("[%s] Already exists in current DB. Skipping...", url.Url)
			continue
		}

		logger.INFO().Printf("[%s] URL does not exist in current DB. Continuing....", url.Url)

		// First add the URL in current DB
		logger.INFO().Printf("[%s] Adding URL...", url.Url)
		if _, err := newDBClient.Queries.AddURL(context.Background(), db.AddURLParams{
			ID:  url.ID,
			Url: url.Url,
		}); err != nil {
			logger.ERROR().Printf("[%s] Failed to add URL!\n%s\nMoving to next one...", url.Url, err.Error())
			continue
		}
		logger.SUCCESS().Printf("[%s] URL added successfully!", url.Url)

		// Get credentials of the url from old DB
		credentials, err := oldDBclient.Queries.GetCredentials(context.Background(), db.GetCredentialsParams{
			Url: url.Url,
		})
		if err != nil {
			logger.ERROR().Printf("[%s] Failed to fetch credentials.\n%s\nSkipping...", url.Url, err.Error())
			continue
		}

		// Copy each credential into current DB
		for i, cred := range credentials {
			logger.INFO().Printf("[%s] Adding credential %d", url.Url, i+1)

			if _, err := newDBClient.Queries.AddCredential(context.Background(), db.AddCredentialParams{
				ID:             cred.ID,
				Url:            url.Url,
				FormInputID:    cred.FormInputID,
				FormInputName:  cred.FormInputName,
				FormInputXpath: cred.FormInputXpath,
				FormInputType:  cred.FormInputType,
				FormInputVal:   cred.FormInputVal,
			}); err != nil {
				logger.ERROR().Printf("[%s] Failed to add credential!\n%s\nMoving to next one...", url.Url, err.Error())
				continue
			}

			logger.SUCCESS().Printf("[%s] Successfully added credential %d", url.Url, i+1)
		}

		logger.SUCCESS().Printf("[%s] Successfully copied to current DB.\n", url.Url)
	}

	logger.SUCCESS().Println("Data migrated successfully!")
	return nil
}
