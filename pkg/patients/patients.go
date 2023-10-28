package patients

import (
	"ktrhportal/database"
	"ktrhportal/models"
	"ktrhportal/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetPatients(c *gin.Context) {
	db := database.DB
	var patients []models.Patient
	if err := db.Preload(clause.Associations).Find(&patients).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "patients", patients)
}

func GetPatient(c *gin.Context) {
	var payload struct {
		SearchParam string `json:"search_param" binding:"required"`
		SearchVal   string `json:"search_val" binding:"required"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	var entities []models.Patient
	db := database.DB
	if payload.SearchParam == "phone_no" {
		db.Where("Phone = ?", payload.SearchVal).Preload(clause.Associations).First(&entities)
	}
	utilities.Show(c, http.StatusOK, "success", entities)
}
