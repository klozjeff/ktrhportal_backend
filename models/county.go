package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type County struct {
	ID          uuid.UUID      `json:"id" gorm:"primary_key"`
	Name        string         `json:"name"`
	Capital     string         `json:"capital"`
	Code        string         `json:"code"`
	SubCounties string         `json:"sub_counties"`
	Slug        string         `json:"slug" binding:"required" gorm:"unique"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (county *County) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	county.ID = uuid
	county.Slug = slug.Make(county.Name)
	return err
}
