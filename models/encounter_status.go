package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type EncounterStatus struct {
	ID          uuid.UUID      `json:"id" gorm:"primary_key" binding:"required"`
	Title       string         `json:"title" binding:"required" gorm:"unique"`
	Description string         `json:"description"`
	Slug        string         `json:"slug" binding:"required" gorm:"unique"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (encounterStatus *EncounterStatus) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	encounterStatus.ID = uuid
	encounterStatus.Slug = slug.Make(encounterStatus.Title)
	return err
}
