package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Donation struct {
	ID                    uuid.UUID      `json:"id" gorm:"primary_key"`
	DonationAmount        string         `json:"donation_amount"`
	DonationCurrency      string         `json:"donation_currency"`
	DonationInitiative    string         `json:"donation_initiative"`
	DonationInitiativeID  string         `json:"donation_initiative_id"`
	FullName              string         `json:"full_name"`
	Email                 string         `json:"email"`
	Phone                 string         `json:"phone"`
	DonationPaymentMethod string         `json:"donation_payment_method"`
	CreatedAt             time.Time      `json:"-"`
	UpdatedAt             time.Time      `json:"-"`
	DeletedAt             gorm.DeletedAt `json:"-"`
}

func (donation *Donation) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	donation.ID = uuid
	return err
}
