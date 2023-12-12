package donations

import (
	"ktrhportal/database"
	"ktrhportal/filters"
	"ktrhportal/models"
	"ktrhportal/utilities"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func AddDonation(c *gin.Context) {
	var payload struct {
		FullName              string `json:"full_name" binding:"required"`
		Email                 string `json:"email" binding:"required"`
		Phone                 string `json:"phone" binding:"required"`
		DonationAmount        string `json:"donation_amount" binding:"required"`
		DonationCurrency      string `json:"donation_currency" binding:"required"`
		DonationInitiative    string `json:"donation_initiative" binding:"required"`
		DonationInitiativeID  string `json:"donation_initiative_id" binding:"required"`
		DonationPaymentMethod string `json:"donation_payment_method" binding:"required"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}

	//database connection
	db := database.DB
	donation := models.Donation{
		FullName:              payload.FullName,
		Email:                 payload.Email,
		Phone:                 payload.Phone,
		DonationAmount:        payload.DonationAmount,
		DonationCurrency:      payload.DonationCurrency,
		DonationInitiative:    payload.DonationInitiative,
		DonationInitiativeID:  payload.DonationInitiativeID,
		DonationPaymentMethod: payload.DonationPaymentMethod,
	}
	if err := db.Create(&donation).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "Donation details added successfully", payload)
}

func AllDonations(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.DonationsFilter)
	db := database.DB
	var entities []models.Donation
	if (filters.DonationsFilter{}) == *input {
		if err := db.
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "donations"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Where("donations.full_name ILIKE ?", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "donations"))
			return
		}
	}
	utilities.Show(c, http.StatusOK, "donations", entities)
}
