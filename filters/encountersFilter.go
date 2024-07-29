package filters

import (
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type EncountersFilter struct {
	StatusId   string `in:"query=status"`
	ProviderId string `in:"query=provider"`
	ClientId   string `in:"query=client"`
	DateRange  string `in:"query=dateRange"`
	Global     string `in:"query=global"`
}

func (ecountersFilter *EncountersFilter) EncounterStatusFilter(db *gorm.DB) *gorm.DB {
	return db.
		Where("status_id = ?", ecountersFilter.StatusId)
}
func (ecountersFilter *EncountersFilter) EncounterProviderFilter(db *gorm.DB) *gorm.DB {
	return db.
		Where("provider_id = ?", ecountersFilter.ProviderId)
}
func (ecountersFilter *EncountersFilter) EncounterClientFilter(db *gorm.DB) *gorm.DB {
	return db.
		Where("client_id = ?", ecountersFilter.ClientId)
}

func (ecountersFilter *EncountersFilter) EncounterBydateRangeFilter(db *gorm.DB) *gorm.DB {
	datesArray := strings.Split(ecountersFilter.DateRange, " ")
	from, error := time.Parse("2006-01-02", datesArray[0])
	log.Printf("Date Format Error: %s", error.Error())
	to, error := time.Parse("2006-01-02", datesArray[1])
	log.Printf("Date Format Error: %s", error.Error())
	return db.
		Where("created_at BETWEEN ? AND ?", from, to)
}
