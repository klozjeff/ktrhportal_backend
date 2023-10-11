package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Specialty struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug" binding:"required" gorm:"unique"`
	Description string    `json:"description"`
	CreatedByID string    `json:"created_by"`
	CreatedBy   *User
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func (specialty *Specialty) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	specialty.ID = uuid
	specialty.Slug = slug.Make(specialty.Name)
	return err
}
