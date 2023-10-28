package routes

import (
	"ktrhportal/middlewares"
	Appointments "ktrhportal/pkg/appointments"
	Patients "ktrhportal/pkg/patients"

	"github.com/gin-gonic/gin"
)

func SetupAppRoutes(appRoute *gin.RouterGroup) {
	settings := appRoute.Group("/app")
	{
		settings.POST("/book-appointment", Appointments.AddAppointment)
		settings.GET("/appointments", middlewares.AuthMiddleware(), Appointments.GetAppointments)

		//patients
		settings.GET("/patients", middlewares.AuthMiddleware(), Patients.GetPatients)
		settings.POST("/patient", Patients.GetPatient)

	}
}
