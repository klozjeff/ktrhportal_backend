package routes

import (
	Auth "ktrhportal/auth"
	"ktrhportal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(appRoute *gin.RouterGroup) {
	auth := appRoute.Group("/auth")
	{
		auth.POST("/login", Auth.Login)
		auth.POST("/register", Auth.Register)
		auth.GET("/currentuser", middlewares.AuthMiddleware(), Auth.CurrentUser)
	}
}
