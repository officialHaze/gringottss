package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/officialhaze/gringottss/api-server/helpers"
	"github.com/officialhaze/gringottss/api-server/logger"
)

// Handle building the DB
func HandleBuildingDB(c *gin.Context) {
	_, err := helpers.BuildDB()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, map[string]string{
			"Error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, map[string]string{
		"message": "Success!",
	})
}

// Handle migrating old DB tables to current DB
func HandleOldDBMigration(c *gin.Context) {
	dbname := c.Param("dbname")

	if err := helpers.MigrateOldDB(dbname); err != nil {
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
