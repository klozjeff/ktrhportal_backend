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
		auth.GET("/currentuser", middlewares.AuthMiddleware(), Auth.UserAccount)
		auth.POST("/logout", middlewares.AuthMiddleware(), Auth.Logout)

	}
	account := appRoute.Group("/account")
	{
		account.GET("/profile", middlewares.AuthMiddleware(), Auth.UserAccount)
		account.POST("/update-photo", middlewares.AuthMiddleware(), Auth.UpdateUserProfile)
		account.POST("/remove-photo", middlewares.AuthMiddleware(), Auth.RemoveUserProfile)
	}
	appRoute.GET("/uploads/:filename", Auth.ServeUploadedFile)
}
