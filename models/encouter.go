package models

import (
	"ktrhportal/utilities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Encounter struct {
	ID                 uuid.UUID        `json:"id" gorm:"primary_key"`
	EncounterNumber    string           `json:"number"`
	ClientId           string           `json:"-"`
	Client             *Client          `json:"client"`
	ProviderId         *string          `json:"provider_id"`
	Provider           *Provider        `json:"provider"`
	AppointmentId      *string          `json:"appointment_id"`
	Appointment        *Appointment     `json:"appointment"`
	EncounterStartTime string           `json:"encounter_start_time"`
	EncounterEndTime   *string          `json:"encounter_end_time"`
	EncounterStartDate string           `json:"encounter_start_date"`
	EncounterEndDate   *string          `json:"encounter_end_date"`
	StatusId           string           `json:"-"`
	Status             *EncounterStatus `json:"encounter_status"`
	CreatedBy          string           `json:"created_by"`
	Notes              *[]Note          `json:"-" gorm:"foreignKey:EncounterID"`
	CreatedAt          time.Time        `json:"-"`
	UpdatedAt          time.Time        `json:"-"`
	DeletedAt          gorm.DeletedAt   `json:"-"`
}

func (encounter *Encounter) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	encounter.ID = uuid
	var status string
	scope.Model(EncounterStatus{}).Where("slug = ?", "new").Select("id").First(&status)
	encounter.StatusId = status

	// Get today's date
	today := time.Now().UTC().Format("2006-01-02")

	var count int64
	scope.Model(Encounter{}).Where("DATE(created_at) = ?", today).Count(&count)
	number := utilities.GenerateAutoIncrementNumber(int(count) + 1)
	encounter.EncounterNumber = number

	return err
}

type Note struct {
	ID            uuid.UUID      `json:"id" gorm:"primary_key"`
	EncounterID   *uuid.UUID     `json:"encounter_id,omitempty"`
	AppointmentID *uuid.UUID     `json:"appointment_id,omitempty"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	CreatedBy     string         `json:"created_by"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `json:"-"`
}

func (note *Note) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	note.ID = uuid
	return err
}
