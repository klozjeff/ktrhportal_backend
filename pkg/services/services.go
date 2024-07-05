package services

import (
	"ktrhportal/database"
	"ktrhportal/filters"
	"ktrhportal/middlewares"
	"ktrhportal/models"
	"ktrhportal/services"
	"ktrhportal/utilities"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func AddService(c *gin.Context) {
	var payload struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	//database connection
	db := database.DB
	service := models.Service{
		Name:        payload.Name,
		Description: payload.Description,
		CreatedByID: middlewares.GetAuthUserID(c),
	}
	if err := db.Create(&service).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "Service added successfully", payload)
}

func ListServices(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.ServicesFilter)
	db := database.DB
	var entities []models.Service
	if (filters.ServicesFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "services"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Where("services.name ILIKE ?", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "services"))
			return
		}
	}
	services.PaginationResponse(db, c, http.StatusOK, "services", entities, models.Service{})

}

func GetServiceDetails(c *gin.Context) {
	db := database.DB
	var service models.Service
	if err := db.
		Where("id = ?", c.Param("id")).
		Preload(clause.Associations).
		First(&service).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "success", service)
}

func DeleteService(c *gin.Context) {
	db := database.DB
	var service models.Service
	if err := db.
		Where("id = ?", c.Param("id")).
		First(&service).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.Delete(&service).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Service deleted successfully")
}
