package routes

import (
	"ktrhportal/middlewares"
	Appointments "ktrhportal/pkg/appointments"

	"github.com/gin-gonic/gin"
)

func SetupAppRoutes(appRoute *gin.RouterGroup) {
	settings := appRoute.Group("/app")
	{
		settings.POST("/book-appointment", Appointments.AddAppointment)
		settings.GET("/appointments", middlewares.AuthMiddleware(), Appointments.GetAppointments)

	}
}
