package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PatientFeedback struct {
	ID                          uuid.UUID      `json:"id" gorm:"primary_key"`
	FullName                    string         `json:"first_name"`
	VisitedPlaces               string         `json:"visited_places"`
	VisitReasons                string         `json:"visit_reasons"`
	EasyToMakeAppointment       string         `json:"easy_to_make_appointment"`
	AppointmentshortWaitingTime string         `json:"appointment_short_waiting_time"`
	PoliteStaff                 string         `json:"polite_staff"`
	ListeningDoctors            string         `json:"listening_doctors"`
	ClearExplanations           string         `json:"clear_explanations_by_staff"`
	CaringNurses                string         `json:"caring_nurses"`
	HelpfulBillingStaff         string         `json:"helpful_billing_staff"`
	ConvinientOperationsHours   string         `json:"convinient_ops_hours"`
	CleanFacility               string         `json:"clean_facility"`
	EasyDirections              string         `json:"easy_directions_and_signages"`
	ServicesRating              string         `json:"services_rating"`
	CanRecommendUs              string         `json:"can_recommend_us"`
	Slug                        string         `json:"slug" binding:"required" gorm:"unique"`
	CreatedAt                   time.Time      `json:"-"`
	UpdatedAt                   time.Time      `json:"-"`
	DeletedAt                   gorm.DeletedAt `json:"-"`
}

func (feedback *PatientFeedback) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	feedback.ID = uuid
	return err
}
