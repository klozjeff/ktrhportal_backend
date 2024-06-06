package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Country struct {
	ID                uuid.UUID      `json:"id"`
	Capital           string         `json:"capital"`
	Continent         string         `json:"continent"`
	CountryName       string         `json:"country_name"`
	Slug              string         `json:"slug"`
	Currency          string         `json:"currency"`
	E164              string         `json:"e164"`
	IsoNumeric        string         `json:"iso_numeric"`
	Iso2              string         `json:"iso2"`
	LanguageCodes     string         `json:"language_codes"`
	Languages         string         `json:"languages"`
	TimezoneInCapital string         `json:"time_zone_in_capital"`
	TopLevelDomain    string         `json:"top_level_domain"`
	PhoneCode         string         `json:"phone_code"`
	Fips              string         `json:"fips"`
	Iso3              string         `json:"iso3"`
	Status            string         `json:"status" gorm:"default:inactive"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (country *Country) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	country.ID = uuid
	if country.CountryName == "Kenya" {
		country.Status = "active"
	}
	country.Slug = slug.Make(country.CountryName)
	return err
}
