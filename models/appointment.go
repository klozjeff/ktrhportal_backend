package models

import (
	"ktrhportal/utilities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
type Appointment struct {
	ID                  uuid.UUID `json:"id" gorm:"primary_key"`
	PatientID           string    `json:"patient_id"`
	Patient             *Patient  `json:"patient"`
	DateOfAppointment   string    `json:"doa"`
	TimeOfAppointment   string    `json:"toa"`
	SpecialtyID         string    `json:"specialty"`
	Specialty           *Specialty
	DoctorID            string `json:"doctor_id"`
	Doctor              *Doctor
	SeekingCareFor      string                    `json:"seeking_care_for"`
	RelationshipID      *string                   `json:"relationship_id"`
	Relationship        *Relationship             `json:"relationship"`
	AppointmentStatusID string                    `json:"appointment_status_id"`
	AppointmentStatus   *AppointmentStatus        `json:"status"`
	PaymentMethodID     string                    `json:"payment_method_id"`
	PaymentMethod       *AppointmentPaymentMethod `json:"payment_methods"`
	InsuraceProviderID  *string                   `json:"insurance_provider_id"`
	InsuraceProvider    *InsuranceProvider        `json:"insurance_providers"`
	Slug                string                    `json:"slug" binding:"required" gorm:"unique"`
	CreatedAt           time.Time                 `json:"-"`
	UpdatedAt           time.Time                 `json:"-"`
	DeletedAt           gorm.DeletedAt            `json:"-"`
}

func (appointment *Appointment) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	appointment.ID = uuid
	return err
}
*/

type Appointment struct {
	ID                 uuid.UUID                 `json:"id" gorm:"primary_key"`
	AppointmentNumber  string                    `json:"number" binding:"required" gorm:"unique"`
	ClientID           string                    `json:"-"`
	Client             *Client                   `json:"client"`
	DateOfAppointment  string                    `json:"doa"`
	TimeOfAppointment  string                    `json:"toa"`
	ServiceID          string                    `json:"-"`
	Service            *Service                  `json:"service"`
	ProviderID         string                    `json:"-"`
	Provider           *Provider                 `json:"provider"`
	StatusID           string                    `json:"-"`
	Status             *AppointmentStatus        `json:"status"`
	PaymentMethodID    string                    `json:"-"`
	PaymentMethod      *AppointmentPaymentMethod `json:"payment_method"`
	InsuraceProviderID *string                   `json:"-"`
	InsuraceProvider   *InsuranceProvider        `json:"insurance_provider"`
	CreatedBy          string                    `json:"created_by"`
	Notes              *[]Note                   `json:"-" gorm:"foreignKey:AppointmentID"`
	// New fields for rescheduling tracking
	OriginalDateOfAppointment string    `json:"original_doa"`
	OriginalTimeOfAppointment string    `json:"original_toa"`
	RescheduleCount           int       `json:"reschedule_count"`
	LastRescheduledAt         time.Time `json:"last_rescheduled_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (appointment *Appointment) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	appointment.ID = uuid

	var status string
	scope.Model(AppointmentStatus{}).Where("slug = ?", "requested").Select("id").First(&status)
	appointment.StatusID = status
	// Get today's date
	today := time.Now().UTC().Format("2006-01-02")

	var count int64
	scope.Model(Appointment{}).Where("DATE(created_at) = ?", today).Count(&count)
	number := utilities.GenerateAutoIncrementNumber(int(count) + 1)
	appointment.AppointmentNumber = number

	// Initialize rescheduling fields
	appointment.OriginalDateOfAppointment = appointment.DateOfAppointment
	appointment.OriginalTimeOfAppointment = appointment.TimeOfAppointment
	appointment.RescheduleCount = 0

	return err
}

// RescheduleAppointment updates the appointment with a new date and time
func (appointment *Appointment) RescheduleAppointment(db *gorm.DB, newDate string, newTime string) error {
	var status string
	db.Model(AppointmentStatus{}).Where("slug = ?", "rescheduled").Select("id").First(&status)
	appointment.StatusID = status
	appointment.DateOfAppointment = newDate
	appointment.TimeOfAppointment = newTime
	appointment.RescheduleCount++
	appointment.LastRescheduledAt = time.Now()
	return db.Save(appointment).Error
}
