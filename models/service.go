package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type ServiceStatus int

const (
	StatusInactive ServiceStatus = iota
	StatusActive
)

type Service struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug" binding:"required" gorm:"unique"`
	Description string    `json:"description"`
	CreatedByID string    `json:"created_by"`
	Status      ServiceStatus
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

const DefaultServiceStatus = StatusActive

func (service *Service) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	service.ID = uuid
	service.Slug = slug.Make(service.Name)
	service.Status = DefaultServiceStatus
	return err
}
