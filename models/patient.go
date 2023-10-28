package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Patient struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key"`
	FirstName   string    `json:"first_name"`
	MiddleName  string    `json:"middle_name"`
	LastName    string    `json:"last_name"`
	Phone       string    `json:"phone_no"`
	Email       string    `json:"email_address"`
	Gender      string    `json:"gender"`
	LanguageID  string    `json:"language_id"`
	Language    *Language `json:"language"`
	Address     string    `json:"physical_address"`
	CountyID    string    `json:"county"`
	County      *County
	SubCounty   string         `json:"sub_county"`
	Appointment *[]Appointment `json:"patient_appointments"`
	//CreatedByID string         `json:"created_by"`
	//CreatedBy   *User
	Slug      string         `json:"slug" binding:"required" gorm:"unique"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (patient *Patient) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	patient.ID = uuid
	patient.Slug = slug.Make(patient.FirstName + "_" + patient.LastName)
	return err
}
