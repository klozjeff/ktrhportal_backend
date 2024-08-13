package encounters

import (
	"ktrhportal/database"
	"ktrhportal/filters"
	"ktrhportal/middlewares"
	"ktrhportal/models"
	"ktrhportal/services"
	"ktrhportal/utilities"
	"net/http"
	"time"

	"github.com/ggicci/httpin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	parsedTime, err := time.Parse("2006-01-02", payload.EncounterStartDate)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid date format"})
		return
	}
	formattedDate := parsedTime.Format("2006-01-02")

	encounter := models.Encounter{
		EncounterStartTime: payload.EncounterStartTime,
		EncounterStartDate: formattedDate,
		ClientId:           payload.ClientId,
		CreatedBy:          middlewares.GetAuthUserID(c),
	}

	// Handle ProviderId
	if payload.ProviderId != "" {
		encounter.ProviderId = &payload.ProviderId
	}

	// Handle AppointmentId
	if payload.AppointmentId != "" {
		encounter.AppointmentId = &payload.AppointmentId
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

		if input.StatusId != "" {
			scopes = append(scopes, input.EncounterStatusFilter)
		}
		if input.ProviderId != "" {
			scopes = append(scopes, input.EncounterProviderFilter)
		}
		if input.ClientId != "" {
			scopes = append(scopes, input.EncounterClientFilter)
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

func DeleteEncounter(c *gin.Context) {
	db := database.DB
	var encounter models.Encounter
	if err := db.
		Where("id = ?", c.Param("id")).
		First(&encounter).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.Delete(&encounter).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Encounter deleted successfully")
}

// AddNoteToEncounter adds a new note to an existing encounter
func AddNoteToEncounter(c *gin.Context) {
	var payload struct {
		Title       string    `json:"title" binding:"required"`
		Content     string    `json:"content" binding:"required"`
		EncounterId uuid.UUID `json:"encounter_id" binding:"required"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	db := database.DB
	note := models.Note{
		Title:       payload.Title,
		Content:     payload.Content,
		EncounterID: &payload.EncounterId,
		CreatedBy:   middlewares.GetAuthUserID(c),
	}
	if err := db.Save(&note).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Encounter note created successfully", map[string]interface{}{
		"id": note.ID,
	})
}

func ListEncounterNotes(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.NotesFilter)
	db := database.DB
	var entities []models.Note
	if (filters.NotesFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Where("encounter_id = ?", c.Param("id")).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "notes"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Where("encounter_id = ?", c.Param("id")).
			Where("notes.title ILIKE ? OR notes.content ILIKE ?", "%"+input.Global+"%", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "notes"))
			return
		}
	}
	services.PaginationResponse(db, c, http.StatusOK, "notes", entities, models.Note{})
}

func DeleteNote(c *gin.Context) {
	db := database.DB
	var note models.Note
	if err := db.
		Where("id = ?", c.Param("id")).
		First(&note).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.Delete(&note).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.ShowMessage(c, http.StatusOK, "Note deleted successfully")
}

func GetNoteDetails(c *gin.Context) {
	db := database.DB
	var note models.Note
	if err := db.
		Where("id = ?", c.Param("id")).
		Preload(clause.Associations).
		First(&note).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "success", note)
}
