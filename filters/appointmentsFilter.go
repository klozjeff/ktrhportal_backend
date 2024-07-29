package filters

import (
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AppointmentsFilter struct {
	StatusId   string `in:"query=status"`
	ProviderId string `in:"query=provider"`
	ClientId   string `in:"query=client"`
	DateRange  string `in:"query=dateRange"`
	Global     string `in:"query=global"`
}

func (appointmentsFilter *AppointmentsFilter) AppointmentStatusFilter(db *gorm.DB) *gorm.DB {
	return db.
		Where("status_id = ?", appointmentsFilter.StatusId)
}
func (appointmentsFilter *AppointmentsFilter) AppointmentDoctorFilter(db *gorm.DB) *gorm.DB {
	return db.
		Where("provider_id = ?", appointmentsFilter.ProviderId)
}
func (appointmentsFilter *AppointmentsFilter) AppointmentClientFilter(db *gorm.DB) *gorm.DB {
	return db.
		Where("client_id = ?", appointmentsFilter.ClientId)
}

func (appointmentsFilter *AppointmentsFilter) AppointmentBydateRangeFilter(db *gorm.DB) *gorm.DB {

	datesArray := strings.Split(appointmentsFilter.DateRange, " ")
	if len(datesArray) != 2 {
		log.Printf("Invalid date range format")
		return db
	}

	from, err := time.Parse("2006-01-02", datesArray[0])
	if err != nil {
		log.Printf("Error parsing 'from' date: %s", err.Error())
		return db
	}
	to, err := time.Parse("2006-01-02", datesArray[1])
	if err != nil {
		log.Printf("Error parsing 'to' date: %s", err.Error())
		return db
	}
	return db.Where("date_of_appointment BETWEEN ? AND ?", from.Format("2006-01-02"), to.Format("2006-01-02"))
}
