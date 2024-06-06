package patients

import (
	"ktrhportal/database"
	"ktrhportal/filters"
	"ktrhportal/models"
	"ktrhportal/services"
	"ktrhportal/utilities"
	"net/http"

	"github.com/ggicci/httpin"
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

func AllPatients(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.PatientsFilter)
	db := database.DB
	var entities []models.Patient
	if (filters.PatientsFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "patients"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Joins("LEFT JOIN genders ON genders.id=patients.gender_id").
			Joins("LEFT JOIN languages ON languages.id=patients.language_id").
			Joins("LEFT JOIN counties ON counties.id=patients.county_id").
			Joins("LEFT JOIN sub_counties ON sub_counties.id=patients.sub_county_id").
			Where("patients.first_name ILIKE ? OR patients.middle_name ILIKE ? OR patients.last_name ILIKE ? OR patients.email ILIKE ? OR patients.phone ILIKE ? OR genders.title ILIKE ? OR languages.title ILIKE ? OR counties.name ILIKE ? OR sub_counties.name ILIKE ?", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "patients"))
			return
		}
	}
	//utilities.Show(c, http.StatusOK, "patients", entities)
	services.PaginationResponse(db, c, http.StatusOK, "patients", entities, models.Patient{})

}

func AddPatient(c *gin.Context) {
	var payload struct {
		FirstName  string `json:"first_name" binding:"required"`
		MiddleName string `json:"middle_name"`
		LastName   string `json:"last_name" binding:"required"`
		Phone      string `json:"phone" binding:"required"`
		Email      string `json:"email" binding:"required"`
		Gender     string `json:"gender" binding:"required"`
		Language   string `json:"language" binding:"required"`
		Address    string `json:"physical_address"`
		County     string `json:"county" binding:"required"`
		SubCounty  string `json:"sub_county" binding:"required"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	db := database.DB
	patient := models.Patient{
		FirstName:   payload.FirstName,
		MiddleName:  payload.MiddleName,
		LastName:    payload.LastName,
		Phone:       payload.Phone,
		Email:       payload.Email,
		GenderID:    payload.Gender,
		LanguageID:  payload.Language,
		Address:     payload.Address,
		CountyID:    payload.County,
		SubCountyID: payload.SubCounty,
	}
	if err := db.Create(&patient).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "Patient information added successfully", patient)
}
