package routes

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/officialhaze/gringottss/api-server/api/REST/server/routes/v1"
)

func MapRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// API v1 route controller map
	v1.RouteControllerMap(api)
}
