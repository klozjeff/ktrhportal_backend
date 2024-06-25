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
	Client             *Patient         `json:"client"`
	ProviderId         *string          `json:"-"`
	Provider           *Doctor          `json:"provider"`
	AppointmentId      *string          `json:"-"`
	Appointment        *Appointment     `json:"appointment"`
	EncounterStartTime string           `json:"encounter_start_time"`
	EncounterEndTime   *string          `json:"encounter_end_time"`
	EncounterStartDate string           `json:"encounter_start_date"`
	EncounterEndDate   *string          `json:"encounter_end_date"`
	StatusId           string           `json:"-"`
	Status             *EncounterStatus `json:"encounter_status"`
	CreatedBy          string           `json:"created_by"`
	Notes              *[]Note          `json:"notes" gorm:"many2many:encounter_notes"`
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
	ID          uuid.UUID  `json:"id" gorm:"primary_key"`
	EncounterID string     `json:"-"`
	Encounter   *Encounter `json:"-"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	CreatedBy   string     `json:"created_by"`
}

func (note *Note) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	note.ID = uuid
	return err
}
