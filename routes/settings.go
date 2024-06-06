package routes

import (
	"ktrhportal/filters"
	"ktrhportal/middlewares"
	Settings "ktrhportal/pkg/settings"

	"github.com/gin-gonic/gin"
)

func SetupSettingRoutes(appRoute *gin.RouterGroup) {
	settings := appRoute.Group("/settings")
	{
		settings.POST("/add_specialty", middlewares.AuthMiddleware(), Settings.AddSpecialty)
		settings.GET("/specialities", Settings.GetSpecialities)
		settings.GET("/all_specialties", middlewares.BindInput(filters.SpecialtiesFilter{}), Settings.AllSpecialties)

		settings.POST("/doctor", middlewares.AuthMiddleware(), Settings.AddDoctor)
		settings.GET("/doctors", Settings.GetDoctors)
		settings.POST("/doctor_details", Settings.GetDoctorDetails)
		settings.GET("/doctors/:specialty_id", Settings.GetDoctorBySpecialty)
		settings.POST("/specialty_doctors", Settings.GetDoctorsBySpecialty)
		settings.GET("/doctor/:id", Settings.GetDoctor)
		settings.GET("/all_doctors", middlewares.BindInput(filters.DoctorsFilter{}), Settings.AllDoctors)

		settings.GET("/roles", Settings.GetRoles)
		settings.GET("/relationships", Settings.GetRelationships)
		settings.GET("/genders", Settings.GetGenders)
		settings.GET("/languages", Settings.GetLanguages)
		settings.GET("/insurance_providers", Settings.GetInsuranceProviders)
		settings.GET("/payment_methods", Settings.GetPaymentMethods)
		settings.GET("/countries", middlewares.BindInput(filters.CountriesFilter{}), Settings.GetCountries)
		settings.GET("/counties", Settings.GetCounties)
		settings.GET("/sub_counties/:county_slug", Settings.GetSubCounties)
		settings.GET("/appointment_statuses", Settings.GetAppointmentStatuses)

	}
}
