package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type SubCounty struct {
	ID        uuid.UUID      `json:"id" gorm:"primary_key"`
	Name      string         `json:"name"`
	County    string         `json:"county"`
	Slug      string         `json:"slug" binding:"required" gorm:"unique"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (subCounty *SubCounty) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	subCounty.ID = uuid
	subCounty.Slug = slug.Make(subCounty.Name)
	return err
}
