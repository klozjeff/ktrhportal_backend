package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JSON map[string]interface{}

type FormTemplate struct {
	ID        uuid.UUID      `json:"id" gorm:"primary_key"`
	Name      string         `json:"name"`
	Fields    string         `json:"fields"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (template *FormTemplate) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	template.ID = uuid
	return err
}

type FormSubmission struct {
	ID         uuid.UUID      `json:"id" gorm:"primary_key"`
	TemplateID uuid.UUID      `json:"-"`
	Template   *FormTemplate  `json:"template"`
	Data       string         `json:"data"` // or `gorm:"type:json"` depending on your database
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-"`
}

func (submission *FormSubmission) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	submission.ID = uuid
	return err
}
