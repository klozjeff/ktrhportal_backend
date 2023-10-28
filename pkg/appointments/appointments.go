package appointments

import (
	"errors"
	"fmt"
	"ktrhportal/database"
	"ktrhportal/models"
	"ktrhportal/services"
	"ktrhportal/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func AddAppointment(c *gin.Context) {
	var payload struct {
		FirstName         string `json:"first_name" binding:"required"`
		MiddleName        string `json:"middle_name"`
		LastName          string `json:"last_name" binding:"required"`
		Phone             string `json:"phone" binding:"required"`
		Email             string `json:"email" binding:"required"`
		Gender            string `json:"gender" binding:"required"`
		Language          string `json:"language" binding:"required"`
		Address           string `json:"physical_address"`
		County            string `json:"county" binding:"required"`
		SubCounty         string `json:"sub_county" binding:"required"`
		DateOfAppointment string `json:"appointment_date" binding:"required"`
		TimeOfAppointment string `json:"appointment_time" binding:"required"`
		Specialty         string `json:"specialty" binding:"required"`
		Doctor            string `json:"doctor" binding:"required"`
		PaymentMethod     string `json:"payment_method" binding:"required"`
		InsuranceProvider string `json:"insurance_provider"`
		SeekingCareFor    string `json:"seeking_care_for" binding:"required"`
		Relationship      string `json:"relationship"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}

	db := database.DB
	appointmentCode, err := utilities.GenerateOTP(4)
	if err != nil {
		return
	}

	patient := models.Patient{
		FirstName:  payload.FirstName,
		MiddleName: payload.MiddleName,
		LastName:   payload.LastName,
		Phone:      payload.Phone,
		Email:      payload.Email,
		Gender:     GetEntityIDBySlug(models.Gender{}, payload.Gender),
		LanguageID: GetEntityIDBySlug(models.Language{}, payload.Language),
		Address:    payload.Address,
		CountyID:   GetEntityIDBySlug(models.County{}, payload.County),
		SubCounty:  GetEntityIDBySlug(models.SubCounty{}, payload.SubCounty),
	}

	appointment := models.Appointment{
		DateOfAppointment:   payload.DateOfAppointment,
		TimeOfAppointment:   payload.TimeOfAppointment,
		SpecialtyID:         GetEntityIDBySlug(models.Specialty{}, payload.Specialty),
		DoctorID:            GetEntityIDBySlug(models.Doctor{}, payload.Doctor),
		SeekingCareFor:      payload.SeekingCareFor,
		RelationshipID:      GetEntityIDBySlug(models.Relationship{}, payload.Relationship),
		AppointmentStatusID: GetEntityIDBySlug(models.AppointmentStatus{}, "new"),
		InsuraceProviderID:  GetEntityIDBySlug(models.InsuranceProvider{}, payload.InsuranceProvider),
		PaymentMethodID:     GetEntityIDBySlug(models.AppointmentPaymentMethod{}, payload.PaymentMethod),
		Slug:                appointmentCode,
		Patient:             &patient,
	}

	if err := db.Create(&appointment).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	//Send Notifications
	//msg := "Dear " + payload.FirstName + ",your appointment booking has been sent successfully. We will get back as soon as possible for confirmation. Appointment tracking code is " + appointmentCode
	msg := fmt.Sprintf("Dear %s,your appointment booking has been sent successfully. We will get back as soon as possible for confirmation. Appointment tracking code is %s", payload.FirstName, appointmentCode)
	if _, err := SendSMSNotification(payload.Phone, msg); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "patient appointment added successfully", appointment)

}

func GetAppointments(c *gin.Context) {
	db := database.DB
	var appointments []models.Appointment
	if err := db.Preload(clause.Associations).Find(&appointments).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "appointments", appointments)
}

func GetEntityIDBySlug(model interface{}, slug string) string {
	var id string
	database.DB.Model(&model).Where("slug=?", slug).Select("id").First(&id)
	return id
}

func SendSMSNotification(phone_no string, msg string) (bool, error) {
	if _, sendSmsResErr := services.SendSMS(phone_no, msg); sendSmsResErr != nil {
		return false, errors.New(sendSmsResErr.Error())
	}
	return true, nil
}
