package appointments

import (
	"errors"
	"fmt"
	"ktrhportal/database"
	"ktrhportal/filters"
	"ktrhportal/models"
	"ktrhportal/services"
	"ktrhportal/utilities"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	//Check if patient exists
	var patient models.Patient
	if result := db.Where("phone = ?", payload.Phone).First(&patient); result.RowsAffected <= 0 {
		patient.FirstName = payload.FirstName
		patient.MiddleName = payload.MiddleName
		patient.LastName = payload.LastName
		patient.Phone = payload.Phone
		patient.Email = payload.Email
		patient.GenderID = payload.Gender
		patient.LanguageID = payload.Language
		patient.Address = payload.Address
		patient.CountyID = GetEntityIDBySlug(models.County{}, payload.County)
		patient.SubCountyID = payload.SubCounty
		db.Create(&patient)
	}
	appointment := models.Appointment{
		DateOfAppointment:   payload.DateOfAppointment,
		TimeOfAppointment:   payload.TimeOfAppointment,
		SpecialtyID:         payload.Specialty,
		DoctorID:            payload.Doctor,
		SeekingCareFor:      payload.SeekingCareFor,
		RelationshipID:      &payload.Relationship,
		AppointmentStatusID: GetEntityIDBySlug(models.AppointmentStatus{}, "new"),
		InsuraceProviderID:  &payload.InsuranceProvider,
		PaymentMethodID:     GetEntityIDBySlug(models.AppointmentPaymentMethod{}, payload.PaymentMethod),
		Slug:                appointmentCode,
		PatientID:           patient.ID.String(),
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

func AllAppointments(c *gin.Context) {
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
			Joins("LEFT JOIN patients ON patients.id=appointments.patient_id").
			Joins("LEFT JOIN doctors ON doctors.id=appointments.doctor_id").
			Joins("LEFT JOIN appointment_statuses ON appointment_statuses.id=appointments.appointment_status_id").
			Joins("LEFT JOIN specialties ON specialties.id=appointments.specialty_id").
			Where("patients.first_name ILIKE ? OR patients.middle_name ILIKE ? OR patients.last_name ILIKE ? OR patients.email ILIKE ? OR patients.phone ILIKE ? OR appointment_statuses.title ILIKE ? OR doctors.first_name ILIKE ? OR doctors.last_name ILIKE ? OR appointments.slug ILIKE ? OR specialties.name ILIKE ?", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%").
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

func RecordAppointment(c *gin.Context) {
	var payload struct {
		FirstName         string `json:"firstname" binding:"required"`
		LastName          string `json:"lastname" binding:"required"`
		Phone             string `json:"phone" binding:"required"`
		Email             string `json:"emailaddress" binding:"required"`
		Gender            string `json:"gender" binding:"required"`
		Address           string `json:"address"`
		DateOfAppointment string `json:"doa" binding:"required"`
		TimeOfAppointment string `json:"toa" binding:"required"`
		Specialty         string `json:"specialty" binding:"required"`
		Doctor            string `json:"doctor" binding:"required"`
		PaymentMethod     string `json:"payment_method" binding:"required"`
		InsuranceProvider string `json:"insurance_provider"`
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
	//Check if patient exists
	var patient models.Patient
	if result := db.Where("phone = ?", payload.Phone).First(&patient); result.RowsAffected <= 0 {
		patient.FirstName = payload.FirstName
		patient.MiddleName = ""
		patient.LastName = payload.LastName
		patient.Phone = payload.Phone
		patient.Email = payload.Email
		patient.GenderID = payload.Gender
		patient.LanguageID = ""
		patient.Address = payload.Address
		patient.CountyID = ""
		patient.SubCountyID = ""
		db.Create(&patient)
	}
	var relation = "2e0f436f-8f36-45e2-ad1e-2ae56b08d316"
	appointment := models.Appointment{
		DateOfAppointment:   payload.DateOfAppointment,
		TimeOfAppointment:   payload.TimeOfAppointment,
		SpecialtyID:         payload.Specialty,
		DoctorID:            payload.Doctor,
		SeekingCareFor:      "other",
		RelationshipID:      &relation,
		AppointmentStatusID: GetEntityIDBySlug(models.AppointmentStatus{}, "new"),
		InsuraceProviderID:  &payload.InsuranceProvider,
		PaymentMethodID:     payload.PaymentMethod,
		Slug:                appointmentCode,
		PatientID:           patient.ID.String(),
	}

	if err := db.Create(&appointment).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	//Send Notifications
	//msg := "Dear " + payload.FirstName + ",your appointment booking has been sent successfully. We will get back as soon as possible for confirmation. Appointment tracking code is " + appointmentCode
	//msg := fmt.Sprintf("Dear %s,your appointment booking has been sent successfully. We will get back as soon as possible for confirmation. Appointment tracking code is %s", payload.FirstName, appointmentCode)
	//if _, err := SendSMSNotification(payload.Phone, msg); err != nil {
	//utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
	//return
	//}

	utilities.Show(c, http.StatusOK, "patient appointment added successfully", appointment)

}
