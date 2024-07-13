package appointments

import (
	"errors"
	"fmt"
	"ktrhportal/database"
	"ktrhportal/filters"
	"ktrhportal/middlewares"
	"ktrhportal/models"
	"ktrhportal/services"
	"ktrhportal/utilities"
	"net/http"
	"strings"

	"github.com/ggicci/httpin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddAppointment(c *gin.Context) {
	var payload struct {
		FirstName         string `json:"firstname" binding:"required"`
		LastName          string `json:"lastname" binding:"required"`
		Phone             string `json:"phone" binding:"required"`
		Email             string `json:"emailaddress"`
		Gender            string `json:"gender" binding:"required"`
		Address           string `json:"address" binding:"required"`
		DateOfBirth       string `json:"dob" binding:"required"`
		DateOfAppointment string `json:"doa" binding:"required"`
		TimeOfAppointment string `json:"toa" binding:"required"`
		Service           string `json:"service" binding:"required"`
		Provider          string `json:"provider" binding:"required"`
		PaymentMethod     string `json:"payment_method" binding:"required"`
		InsuranceProvider string `json:"insurance_provider"`
		Note              string `json:"note"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	db := database.DB
	//Check if patient exists
	var client models.Client
	if result := db.Where("phone = ?", payload.Phone).First(&client); result.RowsAffected <= 0 {
		client.FirstName = payload.FirstName
		client.LastName = payload.LastName
		client.Phone = payload.Phone
		client.Email = &payload.Email
		client.GenderID = payload.Gender
		client.Address = &payload.Address
		client.DateOfBirth = payload.DateOfBirth
		client.CreatedBy = middlewares.GetAuthUserID(c)
		db.Create(&client)
	}
	appointment := models.Appointment{
		DateOfAppointment:  payload.DateOfAppointment,
		TimeOfAppointment:  payload.TimeOfAppointment,
		ServiceID:          payload.Service,
		ProviderID:         payload.Provider,
		InsuraceProviderID: &payload.InsuranceProvider,
		PaymentMethodID:    payload.PaymentMethod,
		ClientID:           client.ID.String(),
		CreatedBy:          middlewares.GetAuthUserID(c),
	}

	if err := db.Create(&appointment).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Note != "" {
		note := models.Note{
			AppointmentID: &appointment.ID,
			Title:         `Initial Note`,
			Content:       payload.Note,
			CreatedBy:     middlewares.GetAuthUserID(c),
		}
		if err := db.Create(&note).Error; err != nil {
			utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	phone := strings.TrimPrefix(payload.Phone, "+")
	msg := fmt.Sprintf("Dear %s,your appointment has been scheduled successfully.Tracking code is %s", payload.FirstName, appointment.AppointmentNumber)
	if _, err := SendSMSNotification(phone, msg); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Appointment added successfully", appointment)

}

func SendSMSNotification(phone_no string, msg string) (bool, error) {
	if _, sendSmsResErr := services.SendSMS(phone_no, msg); sendSmsResErr != nil {
		return false, errors.New(sendSmsResErr.Error())
	}
	return true, nil
}

func GetEntityIDBySlug(model interface{}, slug string) string {
	var id string
	database.DB.Model(&model).Where("slug=?", slug).Select("id").First(&id)
	return id
}

func ListAppointments(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.AppointmentsFilter)
	db := database.DB
	var entities []models.Appointment
	var scopes []func(*gorm.DB) *gorm.DB
	if (filters.AppointmentsFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "appointments"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Joins("LEFT JOIN clients ON clients.id=appointments.client_id").
			Joins("LEFT JOIN providers ON providers.id=appointments.provider_id").
			Joins("LEFT JOIN appointment_statuses ON appointment_statuses.id=appointments.status_id").
			Joins("LEFT JOIN services ON services.id=appointments.service_id").
			Where("clients.first_name ILIKE ? OR clients.middle_name ILIKE ? OR clients.last_name ILIKE ? OR clients.email ILIKE ? OR clients.phone ILIKE ? OR appointment_statuses.title ILIKE ? OR providers.first_name ILIKE ? OR providers.last_name ILIKE ? OR appointments.appointment_number ILIKE ? OR services.name ILIKE ?", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "appointments"))
			return
		}
	} else {

		if input.AppointmentStatusID != "" {
			scopes = append(scopes, input.AppointmentStatusFilter)
		}
		if input.AppointmentDoctorID != "" {
			scopes = append(scopes, input.AppointmentDoctorFilter)
		}
		if input.DateRange != "" {
			scopes = append(scopes, input.AppointmentBydateRangeFilter)
		}

		if err := db.
			Scopes(scopes...).
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "appointments"))
			return
		}
	}
	services.PaginationResponse(db, c, http.StatusOK, "appointments", entities, models.Appointment{})

}

func GetAppointmentDetails(c *gin.Context) {
	db := database.DB
	var appointment models.Appointment
	if err := db.
		Where("id = ?", c.Param("id")).
		Preload(clause.Associations).
		First(&appointment).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "success", appointment)
}

func DeleteAppointment(c *gin.Context) {
	db := database.DB
	var appointment models.Appointment
	if err := db.
		Where("id = ?", c.Param("id")).
		First(&appointment).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.Delete(&appointment).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Appointment deleted successfully")
}

func UpdateAppointmentStatus(c *gin.Context) {
	var payload struct {
		AppointmentID string `json:"appointment_id" binding:"required"`
		NewStatusID   string `json:"new_status_id" binding:"required"`
		Comments      string `json:"comments"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, utilities.Validate(validationError))
		return
	}
	db := database.DB
	var appointment models.Appointment
	db.Where("id = ?", payload.AppointmentID).First(&appointment)
	appointment.StatusID = payload.NewStatusID
	if err := db.Save(&appointment).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Comments != "" {
		note := models.Note{
			AppointmentID: &appointment.ID,
			Title:         `Status update`,
			Content:       payload.Comments,
			CreatedBy:     middlewares.GetAuthUserID(c),
		}
		if err := db.Create(&note).Error; err != nil {
			utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	utilities.ShowMessage(c, http.StatusOK, "Appointment status updated successfully")
}

func RescheduleAppointment(c *gin.Context) {
	var payload struct {
		AppointmentID        string `json:"appointment_id" binding:"required"`
		NewDateofAppointment string `json:"doa" binding:"required"`
		NewTimeofAppointment string `json:"toa" binding:"required"`
		Reason               string `json:"reason"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, utilities.Validate(validationError))
		return
	}
	db := database.DB
	var appointment models.Appointment
	if err := db.Where("id = ?", payload.AppointmentID).First(&appointment).Error; err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, err.Error())
		return
	}

	if err := appointment.RescheduleAppointment(db, payload.NewDateofAppointment, payload.NewTimeofAppointment); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	if payload.Reason != "" {
		note := models.Note{
			AppointmentID: &appointment.ID,
			Title:         `Appointment Reschedule`,
			Content:       payload.Reason,
			CreatedBy:     middlewares.GetAuthUserID(c),
		}
		if err := db.Create(&note).Error; err != nil {
			utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	utilities.ShowMessage(c, http.StatusOK, "Appointment rescheduled successfully")
}

func AssignProvider(c *gin.Context) {
	var payload struct {
		AppointmentID string `json:"appointment_id" binding:"required"`
		ProviderID    string `json:"provider_id" binding:"required"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, utilities.Validate(validationError))
		return
	}
	db := database.DB
	var appointment models.Appointment
	db.Where("id = ?", payload.AppointmentID).First(&appointment)
	appointment.ProviderID = payload.ProviderID
	if err := db.Save(&appointment).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.ShowMessage(c, http.StatusOK, "Assigned provider to appointment successfully")
}
