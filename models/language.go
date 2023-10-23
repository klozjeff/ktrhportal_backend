package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Language struct {
	ID        uuid.UUID      `json:"id" gorm:"primary_key" binding:"required"`
	Title     string         `json:"title" binding:"required" gorm:"unique"`
	Slug      string         `json:"slug" binding:"required" gorm:"unique"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (language *Language) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	language.ID = uuid
	language.Slug = slug.Make(language.Title)
	return err
}
