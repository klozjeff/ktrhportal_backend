package routes

import "github.com/gin-gonic/gin"

func Setup(router *gin.Engine) {
	appRoute := router.Group("/api/v1")
	SetupAuthRoutes(appRoute)
	SetupSettingRoutes(appRoute)
	SetupAppRoutes(appRoute)
}
