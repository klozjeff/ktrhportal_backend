package encounters

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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddEncounter(c *gin.Context) {
	var payload struct {
		EncounterStartDate string `json:"start_date" binding:"required"`
		EncounterStartTime string `json:"start_time" binding:"required"`
		ClientId           string `json:"client_id" binding:"required"`
		AppointmentId      string `json:"appointment_id"`
		ProviderId         string `json:"provider_id"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	db := database.DB
	encounter := models.Encounter{
		EncounterStartTime: payload.EncounterStartTime,
		EncounterStartDate: payload.EncounterStartDate,
		ClientId:           payload.ClientId,
		ProviderId:         &payload.ProviderId,
		AppointmentId:      &payload.AppointmentId,
		CreatedBy:          middlewares.GetAuthUserID(c),
	}
	if err := db.Create(&encounter).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Encounter created successfully", map[string]interface{}{
		"id":           encounter.ID,
		"encounter_no": encounter.EncounterNumber,
	})

}

func ListEncounters(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.EncountersFilter)
	db := database.DB
	var entities []models.Encounter
	var scopes []func(*gorm.DB) *gorm.DB
	if (filters.EncountersFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "encounters"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Joins("LEFT JOIN patients ON patients.id=encounters.client_id").
			Joins("LEFT JOIN doctors ON doctors.id=encounters.provider_id").
			Joins("LEFT JOIN encounter_statuses ON encounter_statuses.id=status_id").
			Joins("LEFT JOIN appointments ON appointments.id=encounters.appointment_id").
			Where("patients.first_name ILIKE ? OR patients.middle_name ILIKE ? OR patients.last_name ILIKE ? OR patients.email ILIKE ? OR patients.phone ILIKE ? OR encounter_statuses.title ILIKE ? OR doctors.first_name ILIKE ? OR doctors.last_name ILIKE ? ", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "encounters"))
			return
		}
	} else {

		if input.EncounterStatusID != "" {
			scopes = append(scopes, input.EncounterStatusFilter)
		}
		if input.EncounterProviderID != "" {
			scopes = append(scopes, input.EncounterProviderFilter)
		}
		if input.DateRange != "" {
			scopes = append(scopes, input.EncounterBydateRangeFilter)
		}

		if err := db.
			Scopes(scopes...).
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "encounters"))
			return
		}
	}
	services.PaginationResponse(db, c, http.StatusOK, "encounters", entities, models.Encounter{})

}

func GetEncounterDetails(c *gin.Context) {
	db := database.DB
	var encounter models.Encounter
	if err := db.
		Where("id = ?", c.Param("id")).
		Preload(clause.Associations).
		First(&encounter).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "success", encounter)
}
