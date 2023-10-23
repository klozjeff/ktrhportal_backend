package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type InsuraceProvider struct {
	ID               uuid.UUID               `json:"id" gorm:"primary_key" binding:"required"`
	Title            string                  `json:"title" binding:"required" gorm:"unique"`
	Slug             string                  `json:"slug" binding:"required" gorm:"unique"`
	ProviderStatusID string                  `json:"status_id"`
	ProviderStatus   *InsuraceProviderStatus `json:"status"`
	CreatedAt        time.Time               `json:"-"`
	UpdatedAt        time.Time               `json:"-"`
	DeletedAt        gorm.DeletedAt          `json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (insuraceProvider *InsuraceProvider) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	var statusID string
	scope.Model(InsuraceProviderStatus{}).Select("id").Where("slug=?", "active").First(&statusID)
	insuraceProvider.ID = uuid
	insuraceProvider.Slug = slug.Make(insuraceProvider.Title)
	insuraceProvider.ProviderStatusID = statusID
	return err
}
