package providers

import (
	"ktrhportal/models"
	"ktrhportal/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
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
		Services []models.Specialty `json:"specialties"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}

	//database connection
	//db := database.DB

	/*Checking if doctor already exist
	var u models.Provider
	result := db.Where("email = ?", payload.Email).First(&u)
	if result.RowsAffected > 0 {
		utilities.ShowError(c, http.StatusOK, append(utilities.ErrrsList, "Doctor with this email already exist"))
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
	for _, schedulePayload := range payload.Availability {
		schedule := models.Schedule{
			Day:       schedulePayload.Day,
			StartTime: schedulePayload.StartTime,
			EndTime:   schedulePayload.EndTime,
			Active:    schedulePayload.Active,
			CreatedBy: middlewares.GetAuthUserID(c),
		}
		provider.Schedule = append(provider.Schedule, schedule)
	}

	if err := db.Create(&provider).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "Provider details added successfully", payload)
	*/
}
