package providers

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
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func AddProvider(c *gin.Context) {
	var payload struct {
		FirstName    string `json:"firstname" binding:"required"`
		LastName     string `json:"lastname" binding:"required"`
		Email        string `json:"emailaddress" binding:"required"`
		Phone        string `json:"phone" binding:"required"`
		Salutation   string `json:"salutation" binding:"required"`
		Bio          string `json:"bio"`
		Position     string `json:"position"`
		Availability []struct {
			Day       string `json:"day"`
			Active    bool   `json:"active"`
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
		} `json:"availability"`
		Services []struct {
			ServiceID uuid.UUID `json:"id"`
		} `json:"services"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	//database connection
	db := database.DB

	var u models.Provider
	result := db.Where("email = ?", payload.Email).First(&u)
	if result.RowsAffected > 0 {
		utilities.ShowError(c, http.StatusOK, append(utilities.ErrrsList, "Provider with this email already exist"))
		return
	}
	provider := models.Provider{
		Salutation: payload.Salutation,
		FirstName:  payload.FirstName,
		LastName:   payload.LastName,
		Email:      payload.Email,
		Phone:      payload.Phone,
		Bio:        &payload.Bio,
		Position:   &payload.Position,
		CreatedBy:  middlewares.GetAuthUserID(c),
	}
	if err := db.Create(&provider).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	for _, availability := range payload.Availability {
		schedule := models.ProviderSchedules{
			ProviderID: provider.ID,
			Day:        availability.Day,
			Active:     availability.Active,
			StartTime:  availability.StartTime,
			EndTime:    availability.EndTime,
			CreatedBy:  middlewares.GetAuthUserID(c),
		}
		if err := db.Create(&schedule).Error; err != nil {
			utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	for _, service := range payload.Services {
		services := models.ProviderServices{
			ProviderID: provider.ID,
			ServiceID:  service.ServiceID,
		}
		if err := db.Create(&services).Error; err != nil {
			utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	utilities.Show(c, http.StatusOK, "Provider details added successfully", payload)
}

func ListProviders(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.ProvidersFilter)
	db := database.DB
	var entities []models.Provider
	if (filters.ProvidersFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "providers"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Where("providers.first_name ILIKE ? OR providers.last_name ILIKE ? OR providers.email ILIKE ? OR providers.phone ILIKE ? OR providers.salutation ILIKE ? OR providers.position ILIKE ? ", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "providers"))
			return
		}
	}
	services.PaginationResponse(db, c, http.StatusOK, "providers", entities, models.Provider{})
}

func GetProviderDetails(c *gin.Context) {
	db := database.DB
	var provider models.Provider
	if err := db.
		Where("id = ?", c.Param("id")).
		Preload(clause.Associations).
		First(&provider).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "success", provider)
}

func DeleteProvider(c *gin.Context) {
	db := database.DB
	var provider models.Provider
	if err := db.
		Where("id = ?", c.Param("id")).
		First(&provider).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.Delete(&provider).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Provider deleted successfully")
}
