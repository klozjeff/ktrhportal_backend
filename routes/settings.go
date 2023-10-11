package routes

import (
	"ktrhportal/middlewares"
	Settings "ktrhportal/pkg/settings"

	"github.com/gin-gonic/gin"
)

func SetupSettingRoutes(appRoute *gin.RouterGroup) {
	settings := appRoute.Group("/settings")
	{
		settings.POST("/specialty", middlewares.AuthMiddleware(), Settings.AddSpecialty)
		settings.GET("/specialities", Settings.GetSpecialities)
		settings.POST("/doctor", middlewares.AuthMiddleware(), Settings.AddDoctor)
		settings.GET("/doctors", Settings.GetDoctors)
		settings.GET("/doctors/:specialty_id", Settings.GetDoctorBySpecialty)
		settings.GET("/doctor/:id", Settings.GetDoctor)
		settings.GET("/roles", Settings.GetRoles)
		settings.GET("/relationships", Settings.GetRelationships)
		settings.GET("/genders", Settings.GetGenders)
	}
}
