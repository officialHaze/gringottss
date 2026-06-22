package controllers

import (
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
	"github.com/officialhaze/gringottss/api-server/logger"
	"github.com/officialhaze/gringottss/api-server/util"
)

// Handle generating encryption keys (paseto supported)
func HandleEncryptionKeysGen(c *gin.Context) {
	file := path.Join("encryption_keys.yml")

	// Prevent accidental data loss by overwriting existing encryption keys
	if _, err := os.Stat(file); err == nil {
		// Encryption keys already exist. Skip overwrite.
		logger.ERROR().Println("Encryption keys already exist. Skipping overwrite to prevent accidental data loss.")
		c.IndentedJSON(http.StatusMethodNotAllowed, map[string]string{
			"Error": "Encryption keys already exist. Skipping overwrite to prevent accidental data loss.",
		})
		return
	}

	// Generate encryption keys for each type
	encryptionKeys := util.EncryptionKeys{
		PWD_ENCRYPTION_KEY: util.Generate32ByteString(),
	}

	// Save the generated keys in file
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logger.ERROR().Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		return
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	defer enc.Close()
	if err := enc.Encode(encryptionKeys); err != nil {
		logger.ERROR().Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, map[string]string{
		"message": "Success!",
	})
}
