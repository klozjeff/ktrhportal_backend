package routes

import (
	"ktrhportal/filters"
	"ktrhportal/middlewares"
	Appointments "ktrhportal/pkg/appointments"
	Clients "ktrhportal/pkg/clients"
	Donations "ktrhportal/pkg/donations"
	Encounters "ktrhportal/pkg/encounters"
	Feedbacks "ktrhportal/pkg/feedback"
	Patients "ktrhportal/pkg/patients"
	Providers "ktrhportal/pkg/providers"
	Services "ktrhportal/pkg/services"

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

		//Providers
		settings.POST("/add_provider", middlewares.AuthMiddleware(), Providers.AddProvider)
		settings.GET("/providers", middlewares.BindInput(filters.ProvidersFilter{}), Providers.ListProviders)
		settings.GET("/providers/:id", middlewares.AuthMiddleware(), Providers.GetProviderDetails)
		settings.DELETE("/providers/:id", middlewares.AuthMiddleware(), Providers.DeleteProvider)

		//Clients
		settings.POST("/add_client", middlewares.AuthMiddleware(), Clients.AddClient)
		settings.GET("/clients", middlewares.BindInput(filters.ClientsFilter{}), Clients.ListClients)

		//Services
		settings.POST("/add_service", middlewares.AuthMiddleware(), Services.AddService)
		settings.GET("/services", middlewares.BindInput(filters.ServicesFilter{}), Services.ListServices)
		settings.GET("/services/:id", middlewares.AuthMiddleware(), Services.GetServiceDetails)
		settings.DELETE("/services/:id", middlewares.AuthMiddleware(), Services.DeleteService)

	}
}
