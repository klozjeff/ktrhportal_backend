package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Client struct {
	ID          uuid.UUID       `json:"id" gorm:"primary_key"`
	FirstName   string          `json:"first_name"`
	MiddleName  *string         `json:"middle_name"`
	LastName    string          `json:"last_name"`
	Phone       string          `json:"phone_no"`
	Email       *string         `json:"email_address"`
	GenderID    string          `json:"-"`
	Gender      *Gender         `json:"gender"`
	LanguageID  *string         `json:"-"`
	Language    *Language       `json:"language"`
	DateOfBirth string          `json:"dob"`
	Address     *string         `json:"physical_address"`
	CountryID   *string         `json:"-"`
	Country     *Country        `json:"country"`
	City        *string         `json:"city"`
	State       *string         `json:"state"`
	PostalCode  *string         `json:"postalcode"`
	ZipCode     *string         `json:"zipcode"`
	People      *[]ClientPeople `json:"people"`
	Slug        string          `json:"slug" binding:"required"`
	CreatedBy   string          `json:"created_by"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"-"`
	DeletedAt   gorm.DeletedAt  `json:"-"`
}

func (client *Client) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	client.ID = uuid
	client.Slug = slug.Make(client.FirstName + "_" + client.LastName)
	return err
}

type ClientPeople struct {
	ID           uuid.UUID      `json:"id" gorm:"primary_key"`
	ClientID     uuid.UUID      `json:"-"`
	Client       *Client        `json:"-"`
	FullName     string         `json:"full_names"`
	Phone        string         `json:"phone_no"`
	Relationship string         `json:"relationship"`
	CreatedBy    string         `json:"-"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}

func (people *ClientPeople) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	people.ID = uuid
	return err
}
