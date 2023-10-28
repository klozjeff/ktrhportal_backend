package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Appointment struct {
	ID                  uuid.UUID `json:"id" gorm:"primary_key"`
	PatientID           string    `json:"patient_id"`
	Patient             *Patient  `json:"patient"`
	DateOfAppointment   string    `json:"doa"`
	TimeOfAppointment   string    `json:"toa"`
	SpecialtyID         string    `json:"specialty"`
	Specialty           *Specialty
	DoctorID            string `json:"doctor"`
	Doctor              *Doctor
	SeekingCareFor      string                    `json:"seeking_care_for"`
	RelationshipID      string                    `json:"relationship_id"`
	Relationship        *Relationship             `json:"relationship"`
	AppointmentStatusID string                    `json:"account_status_id"`
	AppointmentStatus   *AppointmentStatus        `json:"status"`
	PaymentMethodID     string                    `json:"payment_method_id"`
	PaymentMethod       *AppointmentPaymentMethod `json:"payment_methods"`
	InsuraceProviderID  string                    `json:"insurance_provider_id"`
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
