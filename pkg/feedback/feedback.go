package feedback

import (
	"ktrhportal/database"
	"ktrhportal/models"
	"ktrhportal/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func AddFeedback(c *gin.Context) {
	var payload struct {
		FullName         string `json:"full_name" binding:"required"`
		VisitedPlaces    string `json:"visited_places" binding:"required"`
		VisitReasons     string `json:"visit_reasons" binding:"required"`
		AboutAppointment struct {
			EasyToMakeAppointment       string `json:"easy_to_make" binding:"required"`
			AppointmentshortWaitingTime string `json:"short_waiting_time" binding:"required"`
		} `json:"about_appointment" binding:"required"`
		AboutStaff struct {
			PoliteStaff         string `json:"polite_staff" binding:"required"`
			ListeningDoctors    string `json:"listening_doctors" binding:"required"`
			ClearExplanations   string `json:"clear_explanations" binding:"required"`
			CaringNurses        string `json:"caring_nurses" binding:"required"`
			HelpfulBillingStaff string `json:"helpful_billing_staff" binding:"required"`
		} `json:"about_staff" binding:"required"`
		AboutFacility struct {
			ConvinientOperationsHours string `json:"convinient_ops_hours" binding:"required"`
			CleanFacility             string `json:"clean_facility" binding:"required"`
			EasyDirections            string `json:"easy_directions" binding:"required"`
		} `json:"about_facility" binding:"required"`
		ServicesRating string `json:"services_rating" binding:"required"`
		CanRecommendUs string `json:"can_recommend_us" binding:"required"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	db := database.DB
	feedbackCode, err := utilities.GenerateOTP(4)
	if err != nil {
		return
	}

	feedback := models.PatientFeedback{
		FullName:                    payload.FullName,
		VisitedPlaces:               payload.VisitedPlaces,
		VisitReasons:                payload.VisitReasons,
		EasyToMakeAppointment:       payload.AboutAppointment.EasyToMakeAppointment,
		AppointmentshortWaitingTime: payload.AboutAppointment.EasyToMakeAppointment,
		PoliteStaff:                 payload.AboutStaff.PoliteStaff,
		ListeningDoctors:            payload.AboutStaff.ListeningDoctors,
		ClearExplanations:           payload.AboutStaff.ClearExplanations,
		CaringNurses:                payload.AboutStaff.CaringNurses,
		HelpfulBillingStaff:         payload.AboutStaff.HelpfulBillingStaff,
		ConvinientOperationsHours:   payload.AboutFacility.ConvinientOperationsHours,
		CleanFacility:               payload.AboutFacility.CleanFacility,
		EasyDirections:              payload.AboutFacility.EasyDirections,
		ServicesRating:              payload.ServicesRating,
		CanRecommendUs:              payload.CanRecommendUs,
		Slug:                        feedbackCode,
	}

	if err := db.Create(&feedback).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "feedback added successfully", payload)
}

func GetFeedbacks(c *gin.Context) {
	db := database.DB
	var feedbacks []models.PatientFeedback
	if err := db.Preload(clause.Associations).Find(&feedbacks).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "feedbacks", feedbacks)
}
