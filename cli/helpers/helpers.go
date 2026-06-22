package helpers

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/joho/godotenv"
	"github.com/officialHaze/gringottss/cli/logger"
)

var (
	client = http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

func gringottssEngineAPIBaseUrl() string {
	envs, err := godotenv.Read(path.Join(".env"))
	if err != nil {
		logger.ERROR().Println(err.Error())
		return ""
	}

	host := envs["GRINGOTTSS_ENGINE_HOST"]
	port := envs["GRINGOTTSS_ENGINE_PORT"]

	return fmt.Sprintf("http://%s:%s/api", host, port)
}

// This method calls the Gringottss Engine to build its DB.
func BuildDB(ctx context.Context) {
	logger.INFO().Println("Building Gringottss DB....")

	endpoint := fmt.Sprintf("%s/v1/db/build", gringottssEngineAPIBaseUrl())
	fmt.Printf("Endpoint: %s\n", endpoint)
	req, _ := http.NewRequestWithContext(ctx, "POST", endpoint, nil)

	res, err := client.Do(req)
	if err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
		return
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		logger.ERROR().Println(string(b))
		os.Exit(1)
		return
	}

	logger.SUCCESS().Println("Gringottss DB built successfully.")
}

// This method calls the Gringottss Engine to migrate old DB tables to current DB.
// NOTE:- Make sure the old DB is present inside <PROJECT_ROOT>/engine/migrate. Eg:- <PROJECT_ROOT>/engine/migrate/old.db
func MigrateOldDB(ctx context.Context, dbname string) {
	logger.INFO().Printf("Migrating %s....", dbname)

	endpoint := fmt.Sprintf("%s/v1/db/migrate/%s", gringottssEngineAPIBaseUrl(), url.QueryEscape(dbname))
	fmt.Printf("Endpoint: %s\n", endpoint)
	req, _ := http.NewRequestWithContext(ctx, "POST", endpoint, nil)

	res, err := client.Do(req)
	if err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
		return
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		logger.ERROR().Println(string(b))
		os.Exit(1)
		return
	}

	logger.SUCCESS().Printf("%s migrated successfully.", dbname)
}

// Generate encryption keys
func GenerateEncryptionKeys(ctx context.Context) {
	logger.INFO().Println("Generating encryption keys....")

	endpoint := fmt.Sprintf("%s/v1/keys/encryption/generate", gringottssEngineAPIBaseUrl())
	fmt.Printf("Endpoint: %s\n", endpoint)
	req, _ := http.NewRequestWithContext(ctx, "POST", endpoint, nil)

	res, err := client.Do(req)
	if err != nil {
		logger.ERROR().Println(err.Error())
		os.Exit(1)
		return
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		logger.ERROR().Println(string(b))
		os.Exit(1)
		return
	}

	logger.SUCCESS().Println("Encryption keys generated successfully. Restart the gringottss engine.")
}
