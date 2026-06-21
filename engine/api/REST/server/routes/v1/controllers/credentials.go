package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
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

	// check if credential is a password, then hash else skip
	if params.FormInputType == "password" {
		// encrypt the value
		now := time.Now()

		jsonToken := paseto.JSONToken{
			Issuer:   "gringottss_api",
			Subject:  "pwd_encryption",
			IssuedAt: now,
		}
		jsonToken.Set("pwd", params.FormInputVal)
		footer := fmt.Sprintf("pwd_encryption_%s", params.ID)
		encrypted, err := util.PasetoEncrypt(jsonToken, os.Getenv("PWD_ENCRYPTION_KEY"), footer)
		if err != nil {
			logger.ERROR().Println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, map[string]string{
				"Error": "Internal server error!",
			})
			return
		}

		// Set the encrypted value as the new form input value
		params.FormInputVal = encrypted
	}

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
		// decrypt its value and re-set
		for i, cred := range credentials {
			credptr := &cred

			if cred.FormInputType == "password" {
				// decrypt
				token, _, err := util.PasetoDecrypt(cred.FormInputVal, os.Getenv("PWD_ENCRYPTION_KEY"))
				if err != nil {
					logger.ERROR().Println(err.Error())
					c.IndentedJSON(http.StatusInternalServerError, map[string]string{
						"Error": "Internal server error!",
					})
					return
				}

				// Set the decrypted value
				credptr.FormInputVal = token.Get("pwd")

				// Replace the current credential item with the modified one
				tmp := credentials[i]
				credentials[i] = credentials[len(credentials)-1]
				credentials[len(credentials)-1] = tmp
				credentials = credentials[:len(credentials)-1] // remove the current credential item
				credentials = append(credentials, *credptr)    // replace the modified credential item

				// Maintain the original position by re-positioning
				tmp = credentials[i]
				credentials[i] = credentials[len(credentials)-1]
				credentials[len(credentials)-1] = tmp
			}
		}
	}

	c.IndentedJSON(http.StatusOK, map[string]any{
		"message": "Credentials found.",
		"data":    credentials,
	})
}
