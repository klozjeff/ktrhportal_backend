package routes

import (
	"ktrhportal/filters"
	"ktrhportal/middlewares"
	Appointments "ktrhportal/pkg/appointments"
	Clients "ktrhportal/pkg/clients"
	Dashboard "ktrhportal/pkg/dashboard"
	Donations "ktrhportal/pkg/donations"
	Encounters "ktrhportal/pkg/encounters"
	Feedbacks "ktrhportal/pkg/feedback"
	Forms "ktrhportal/pkg/formsmanager"
	Patients "ktrhportal/pkg/patients"
	Providers "ktrhportal/pkg/providers"
	Services "ktrhportal/pkg/services"

	"github.com/gin-gonic/gin"
)

func SetupAppRoutes(appRoute *gin.RouterGroup) {
	settings := appRoute.Group("/app")
	{
		//settings.POST("/book-appointment", Appointments.AddAppointment)
		settings.POST("/add_appointment", middlewares.AuthMiddleware(), Appointments.AddAppointment)
		//settings.GET("/appointments", middlewares.AuthMiddleware(), Appointments.GetAppointments)
		settings.GET("/appointments", middlewares.BindInput(filters.AppointmentsFilter{}), Appointments.ListAppointments)
		settings.GET("/appointments/:id", middlewares.AuthMiddleware(), Appointments.GetAppointmentDetails)
		settings.POST("/appointments/update-status", middlewares.AuthMiddleware(), Appointments.UpdateAppointmentStatus)
		settings.POST("/appointments/reschedule", middlewares.AuthMiddleware(), Appointments.RescheduleAppointment)
		settings.POST("/appointments/assign-provider", middlewares.AuthMiddleware(), Appointments.AssignProvider)
		settings.DELETE("/appointments/:id", middlewares.AuthMiddleware(), Appointments.DeleteAppointment)

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
		settings.DELETE("/encounters/:id", middlewares.AuthMiddleware(), Encounters.DeleteEncounter)
		settings.POST("/encounters/append-note", middlewares.AuthMiddleware(), Encounters.AddNoteToEncounter)
		settings.GET("/encounters/:id/notes", middlewares.BindInput(filters.NotesFilter{}), Encounters.ListEncounterNotes)
		settings.DELETE("/notes/:id", middlewares.AuthMiddleware(), Encounters.DeleteNote)
		settings.GET("/notes/:id", middlewares.AuthMiddleware(), Encounters.GetNoteDetails)

		//Providers
		settings.POST("/add_provider", middlewares.AuthMiddleware(), Providers.AddProvider)
		settings.GET("/providers", middlewares.BindInput(filters.ProvidersFilter{}), Providers.ListProviders)
		settings.GET("/providers/:id", middlewares.AuthMiddleware(), Providers.GetProviderDetails)
		settings.DELETE("/providers/:id", middlewares.AuthMiddleware(), Providers.DeleteProvider)

		//Clients
		settings.POST("/add_client", middlewares.AuthMiddleware(), Clients.AddClient)
		settings.GET("/clients", middlewares.BindInput(filters.ClientsFilter{}), Clients.ListClients)
		settings.GET("/clients/:id", middlewares.AuthMiddleware(), Clients.GetClientDetails)

		//Services
		settings.POST("/add_service", middlewares.AuthMiddleware(), Services.AddService)
		settings.GET("/services", middlewares.BindInput(filters.ServicesFilter{}), Services.ListServices)
		settings.GET("/services/:id", middlewares.AuthMiddleware(), Services.GetServiceDetails)
		settings.DELETE("/services/:id", middlewares.AuthMiddleware(), Services.DeleteService)

		//Forms
		settings.POST("/add_template", middlewares.AuthMiddleware(), Forms.AddTemplate)
		settings.GET("/form_templates", middlewares.BindInput(filters.FormsManagerFilter{}), Forms.ListTemplates)
		settings.GET("/form_templates/:id", middlewares.AuthMiddleware(), Forms.GetTemplateDetails)
		settings.POST("/update_template", middlewares.AuthMiddleware(), Forms.UpdateTemplate)
		settings.DELETE("/form_templates/:id", middlewares.AuthMiddleware(), Forms.DeleteTemplate)
		settings.POST("/form_submissions", middlewares.AuthMiddleware(), Forms.FormSubmmision)
		settings.GET("/form_submissions", middlewares.AuthMiddleware(), Forms.ListFormSubmissions)
		settings.DELETE("/form_submissions/:id", middlewares.AuthMiddleware(), Forms.DeleteFormSubmission)

		//Dashboard
		settings.GET("/dashboard/overview", middlewares.AuthMiddleware(), Dashboard.Overview)

	}
}
