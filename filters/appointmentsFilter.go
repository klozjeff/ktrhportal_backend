package filters

import (
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AppointmentsFilter struct {
	AppointmentStatusID string `in:"query=status"`
	AppointmentDoctorID string `in:"query=doctor"`
	DateRange           string `in:"query=dateRange"`
	Global              string `in:"query=global"`
}

func (appointmentsFilter *AppointmentsFilter) AppointmentStatusFilter(db *gorm.DB) *gorm.DB {
	return db.
		Where("appointment_status_id = ?", appointmentsFilter.AppointmentStatusID)
}
func (appointmentsFilter *AppointmentsFilter) AppointmentDoctorFilter(db *gorm.DB) *gorm.DB {
	return db.
		Where("doctor_id = ?", appointmentsFilter.AppointmentDoctorID)
}

func (appointmentsFilter *AppointmentsFilter) AppointmentBydateRangeFilter(db *gorm.DB) *gorm.DB {
	datesArray := strings.Split(appointmentsFilter.DateRange, " ")
	from, error := time.Parse("2006-01-02", datesArray[0])
	log.Printf("Date Format Error: %s", error.Error())
	to, error := time.Parse("2006-01-02", datesArray[1])
	log.Printf("Date Format Error: %s", error.Error())
	return db.
		Where("created_at BETWEEN ? AND ?", from, to)
}
