package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/officialhaze/gringottss/api-server/db/sqlc"
	"github.com/officialhaze/gringottss/api-server/logger"
	"github.com/officialhaze/gringottss/api-server/sqliteclient"
)

// Handle getting all urls
func HandleURLsFetch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	urls, err := sqliteclient.Queries.ListURLs(ctx)
	if err != nil {
		logger.ERROR().Println(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, map[string]string{
			"Error": "Internal server error!",
		})
		return
	}

	if len(urls) <= 0 {
		urls = []db.Url{}
	}

	c.IndentedJSON(http.StatusOK, urls)
}

// Handle deleting a url
func HandleURLDelete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := c.Query("url")

	err := sqliteclient.Queries.DeleteURL(ctx, url)
	if err != nil {
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
