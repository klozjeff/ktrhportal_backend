package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Provider struct {
	ID         uuid.UUID            `json:"id" gorm:"primary_key"`
	Salutation string               `json:"prefix"`
	FirstName  string               `json:"first_name"`
	LastName   string               `json:"last_name"`
	Email      string               `json:"email"`
	Phone      string               `json:"phone"`
	Services   *[]ProviderServices  `json:"services" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Schedules  *[]ProviderSchedules `json:"schedules" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Position   *string              `json:"position"`
	Bio        *string              `json:"bio"`
	Slug       string               `json:"slug" binding:"required" gorm:"unique"`
	CreatedBy  string               `json:"created_by"`
	CreatedAt  time.Time            `json:"-"`
	UpdatedAt  time.Time            `json:"-"`
	DeletedAt  gorm.DeletedAt       `json:"-"`
}

func (provider *Provider) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	provider.ID = uuid
	provider.Slug = slug.Make(provider.FirstName + "_" + provider.LastName)
	return err
}

type ProviderSchedules struct {
	ID         uuid.UUID      `json:"id" gorm:"primary_key"`
	ProviderID uuid.UUID      `json:"-"`
	Provider   *Provider      `json:"-"`
	Day        string         `json:"day"`
	StartTime  string         `json:"start_time"`
	EndTime    string         `json:"end_time"`
	Active     bool           `json:"active"`
	CreatedBy  string         `json:"created_by"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-"`
}

func (schedule *ProviderSchedules) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	schedule.ID = uuid
	return err
}

type ProviderServices struct {
	ProviderID uuid.UUID `json:"-"`
	Provider   *Provider `json:"-"`
	ServiceID  uuid.UUID `json:"id"`
	Service    *Service  `json:"-"`
}
