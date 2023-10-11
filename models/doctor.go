package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Doctor struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	SpecialtyID string    `json:"specialty_id"`
	Specialty   *Specialty
	Bio         string `json:"bio"`
	Slug        string `json:"slug" binding:"required" gorm:"unique"`
	CreatedByID string `json:"created_by"`
	CreatedBy   *User
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func (doctor *Doctor) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	doctor.ID = uuid
	doctor.Slug = slug.Make(doctor.FirstName + "_" + doctor.LastName)
	return err
}
