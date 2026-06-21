package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/officialhaze/gringottss/api-server/api/REST/server/routes/v1/controllers"
)

func RouteControllerMap(api *gin.RouterGroup) {
	v1 := api.Group("/v1")

	// ======== Credentials ======== //
	v1.GET("/credentials", controllers.HandleCredentialsFetch)     // Fetch credentials under a url
	v1.POST("/credentials", controllers.HandleCredentialAdd)       // Add credential under a url
	v1.DELETE("/credentials", controllers.HandleCredentialsDelete) // Delete credentials under a url

	// ======== URLs ======== //
	v1.GET("/urls", controllers.HandleURLsFetch)    // List all added URLs
	v1.DELETE("/urls", controllers.HandleURLDelete) // Delete an added URL
}
