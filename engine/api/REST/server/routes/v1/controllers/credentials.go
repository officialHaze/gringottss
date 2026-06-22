package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	db "github.com/officialhaze/gringottss/api-server/db/sqlc"
	"github.com/officialhaze/gringottss/api-server/logger"
	"github.com/officialhaze/gringottss/api-server/sqliteclient"
	"github.com/officialhaze/gringottss/api-server/util"
)

type CredentialUpsertPayload struct {
	URL            string `json:"url"`
	FormInputID    string `json:"formInputId"`
	FormInputType  string `json:"formInputType"`
	FormInputName  string `json:"formInputName"`
	FormInputXPath string `json:"formInputXPath"`
	FormInputVal   string `json:"formInputVal"`
}

// Handle adding credential
func HandleCredentialAdd(c *gin.Context) {
	var payload CredentialUpsertPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.ERROR().Println(err.Error())
		c.IndentedJSON(http.StatusUnprocessableEntity, map[string]string{
			"Error": err.Error(),
		})
		return
	}

	urladdctx, cancelurladdctx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelurladdctx()

	credentialaddctx, cancelcredentialaddctx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelcredentialaddctx()

	// First add the url
	urladdArg := db.AddURLParams{
		ID:  uuid.NewString(),
		Url: payload.URL,
	}

	_, err := sqliteclient.Queries.AddURL(urladdctx, urladdArg)
	if err != nil {
		logger.ERROR().Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, map[string]string{
			"Error": "Internal server error!",
		})
		return
	}

	// Now add the credential
	formInputName := &sql.NullString{
		String: payload.FormInputName,
		Valid:  true,
	}
	if payload.FormInputName == "" {
		formInputName.Valid = false
	}
	params := db.AddCredentialParams{
		ID:             uuid.NewString(),
		Url:            payload.URL,
		FormInputID:    payload.FormInputID,
		FormInputName:  *formInputName,
		FormInputXpath: payload.FormInputXPath,
		FormInputType:  payload.FormInputType,
		FormInputVal:   payload.FormInputVal,
	}

	// encrypt any password input field and set the encrypted value
	encryptAndSetPassword(&params)

	added, err := sqliteclient.Queries.AddCredential(credentialaddctx, params)
	if err != nil {
		logger.ERROR().Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, map[string]string{
			"Error": "Internal server error!",
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, added)
}

// Handle deleting credentials
func HandleCredentialsDelete(c *gin.Context) {
	id := c.Query("id")
	url := c.Query("url")

	deletectx, canceldeletectx := context.WithTimeout(context.Background(), 10*time.Second)
	defer canceldeletectx()

	params := db.DeleteCredentialsParams{
		Url: url,
		ID:  id,
	}

	if err := sqliteclient.Queries.DeleteCredentials(deletectx, params); err != nil {
		logger.ERROR().Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, map[string]string{
			"Error": "Internal server error!",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, map[string]string{
		"message": "Success!",
	})
}

// Handle getting credentials
func HandleCredentialsFetch(c *gin.Context) {
	id := c.Query("id")
	url := c.Query("url")
	_, noPwdEncryption := c.GetQuery("nopwdencryption")

	fetchctx, cancelfetchctx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfetchctx()

	params := db.GetCredentialsParams{
		Url: url,
		ID:  id,
	}

	credentials, err := sqliteclient.Queries.GetCredentials(fetchctx, params)
	if err != nil {
		logger.ERROR().Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, map[string]string{
			"Error": "Internal server error!",
		})
		return
	}

	if len(credentials) <= 0 {
		c.IndentedJSON(http.StatusOK, map[string]any{
			"message": "No credentials found.",
			"data":    nil,
		})
		return
	}

	if noPwdEncryption {
		// Iterate over the credentials and if password is found
		// decrypt its value and set the decrypted value
		for i, cred := range credentials {
			credptr := &cred

			if cred.FormInputType == "password" {
				if err := decryptAndSetPassword(credptr, &credentials, i); err != nil {
					logger.ERROR().Println(err.Error())
					continue
				}
			}
		}
	}

	c.IndentedJSON(http.StatusOK, map[string]any{
		"message": "Credentials found.",
		"data":    credentials,
	})
}

// Encrypt plain text password and set the new encrypted value
func encryptAndSetPassword(params *db.AddCredentialParams) {
	// check if credential is a password, then encrypt the value else skip
	if params.FormInputType == "password" {
		if util.LoadedEncryptionKeys.PWD_ENCRYPTION_KEY == "" {
			// No encryption key set, skip
			logger.WARN().Println("No encryption key set. Skipping...")
			return
		}

		now := time.Now()

		jsonToken := paseto.JSONToken{
			Issuer:   "gringottss_engine",
			Subject:  "pwd_encryption",
			IssuedAt: now,
		}
		jsonToken.Set("pwd", params.FormInputVal)
		footer := fmt.Sprintf("pwd_encryption_%s", params.ID)
		encrypted, err := util.PasetoEncrypt(jsonToken, util.LoadedEncryptionKeys.PWD_ENCRYPTION_KEY, footer)
		if err != nil {
			logger.ERROR().Printf("Password encryption failed! Storing in plain-text.\n%s", err.Error())
			return
		}

		// Set the encrypted value as the new form input value
		params.FormInputVal = encrypted
	}
}

// Decrypt password into plain text and set the decrypted value
func decryptAndSetPassword(credential *db.Credential, credentials *[]db.Credential, idx int) error {
	if util.LoadedEncryptionKeys.PWD_ENCRYPTION_KEY == "" {
		// No encryption key set, skip
		return nil
	}

	// decrypt
	token, _, err := util.PasetoDecrypt(credential.FormInputVal, util.LoadedEncryptionKeys.PWD_ENCRYPTION_KEY)
	if err != nil {
		return err
	}

	// Set the decrypted value
	credential.FormInputVal = token.Get("pwd")

	// Replace the current credential item with the modified one
	tmp := (*credentials)[idx]
	(*credentials)[idx] = (*credentials)[len(*credentials)-1]
	(*credentials)[len(*credentials)-1] = tmp
	*credentials = (*credentials)[:len(*credentials)-1] // remove the current credential item
	*credentials = append(*credentials, *credential)    // replace the modified credential item

	// Maintain the original position by re-positioning
	tmp = (*credentials)[idx]
	(*credentials)[idx] = (*credentials)[len(*credentials)-1]
	(*credentials)[len(*credentials)-1] = tmp

	return nil
}
