package routes

import (
	"ktrhportal/filters"
	"ktrhportal/middlewares"
	Appointments "ktrhportal/pkg/appointments"
	Donations "ktrhportal/pkg/donations"
	Encounters "ktrhportal/pkg/encounters"
	Feedbacks "ktrhportal/pkg/feedback"
	Patients "ktrhportal/pkg/patients"

	"github.com/gin-gonic/gin"
)

func SetupAppRoutes(appRoute *gin.RouterGroup) {
	settings := appRoute.Group("/app")
	{
		settings.POST("/book-appointment", Appointments.AddAppointment)
		settings.POST("/add_appointment", Appointments.RecordAppointment)
		settings.GET("/appointments", middlewares.AuthMiddleware(), Appointments.GetAppointments)
		settings.GET("/all_appointments", middlewares.BindInput(filters.AppointmentsFilter{}), Appointments.AllAppointments)
		settings.GET("/appointments/:id", middlewares.AuthMiddleware(), Appointments.GetAppointmentDetails)

		//patients
		settings.GET("/patients", middlewares.AuthMiddleware(), Patients.GetPatients)
		settings.GET("/all_patients", middlewares.BindInput(filters.PatientsFilter{}), Patients.AllPatients)
		settings.POST("/patient", Patients.GetPatient)
		settings.POST("/add_patient", Patients.AddPatient)

		//Feedbacks
		settings.GET("/feedbacks", middlewares.AuthMiddleware(), Feedbacks.GetFeedbacks)
		settings.POST("/add_feedback", Feedbacks.AddFeedback)
		settings.GET("/all_feedbacks", middlewares.BindInput(filters.FeedbacksFilter{}), Feedbacks.AllFeedbacks)

		//Donations
		settings.POST("/add_donation", Donations.AddDonation)
		settings.GET("/all_donations", middlewares.BindInput(filters.DonationsFilter{}), Donations.AllDonations)

		//Encounters
		settings.POST("/add_encounter", middlewares.AuthMiddleware(), Encounters.AddEncounter)
		settings.GET("/encounters", middlewares.BindInput(filters.EncountersFilter{}), Encounters.ListEncounters)
		settings.GET("/encounters/:id", middlewares.AuthMiddleware(), Encounters.GetEncounterDetails)

	}
}
