package settings

import (
	"ktrhportal/database"
	"ktrhportal/middlewares"
	"ktrhportal/models"
	"ktrhportal/pkg/appointments"
	"ktrhportal/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetGenders(c *gin.Context) {
	db := database.DB
	var entities []models.Gender

	if err := db.Find(&entities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusFound, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "genders", entities)
}

func GetRoles(c *gin.Context) {
	db := database.DB
	var entities []models.Role

	if err := db.Find(&entities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusFound, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "roles", entities)
}

func GetRelationships(c *gin.Context) {
	db := database.DB
	var entities []models.Relationship

	if err := db.Find(&entities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusFound, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "relationships", entities)
}
func GetAccountStatuses(c *gin.Context) {
	db := database.DB
	var entities []models.AccountStatus

	if err := db.Find(&entities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusFound, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "account statuses", entities)
}
func GetLanguages(c *gin.Context) {
	db := database.DB
	var entities []models.Language

	if err := db.Find(&entities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusFound, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "languages", entities)
}
func GetPaymentMethods(c *gin.Context) {
	db := database.DB
	var entities []models.AppointmentPaymentMethod

	if err := db.Find(&entities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusFound, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "payment_methods", entities)
}
func GetInsuranceProviders(c *gin.Context) {
	db := database.DB
	var entities []models.InsuranceProvider

	if err := db.Find(&entities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusFound, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "providers", entities)
}
func GetCounties(c *gin.Context) {
	db := database.DB
	var entities []models.County

	if err := db.Find(&entities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusFound, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "counties", entities)
}
func GetSubCounties(c *gin.Context) {
	db := database.DB
	var entities []models.SubCounty
	if err := db.Where("county=?", c.Param("county_slug")).Preload(clause.Associations).Find(&entities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusFound, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "sub-counties", entities)
}

// Specialities
func AddSpecialty(c *gin.Context) {
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
	specialty := models.Specialty{
		Name:        payload.Name,
		Description: payload.Description,
		CreatedByID: middlewares.GetAuthUserID(c),
	}
	if err := db.Create(&specialty).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "specialty added successfully", payload)
}

func GetSpecialities(c *gin.Context) {
	db := database.DB
	var specialities []models.Specialty
	if err := db.Preload(clause.Associations).Find(&specialities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "specialities", specialities)

}

// Doctors
func AddDoctor(c *gin.Context) {
	var payload struct {
		FirstName   string `json:"first_name" binding:"required"`
		LastName    string `json:"last_name" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Phone       string `json:"phone" binding:"required"`
		SpecialtyID string `json:"specialty" binding:"required"`
		Bio         string `json:"bio"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}

	//database connection
	db := database.DB

	//Checking if doctor already exist
	var u models.Doctor
	result := db.Where("email = ?", payload.Email).First(&u)
	if result.RowsAffected > 0 {
		utilities.ShowError(c, http.StatusOK, append(utilities.ErrrsList, "Doctor with this email already exist"))
		return
	}
	// create
	doctor := models.Doctor{
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		Phone:       payload.Phone,
		SpecialtyID: payload.SpecialtyID,
		Bio:         payload.Bio,
		CreatedByID: middlewares.GetAuthUserID(c),
	}
	if err := db.Create(&doctor).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "doctor details added successfully", payload)
}

func GetDoctors(c *gin.Context) {
	db := database.DB
	var doctors []models.Doctor
	if err := db.Preload(clause.Associations).Find(&doctors).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "success", doctors)
}

func GetDoctor(c *gin.Context) {
	db := database.DB
	var doctor models.Doctor
	if err := db.Where("id=?", c.Param("id")).Preload(clause.Associations).Preload("Specialty").Find(&doctor).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "success", doctor)
}
func GetDoctorBySpecialty(c *gin.Context) {
	var entities []models.Doctor
	db := database.DB
	db.Where("specialty_id = ?", c.Param("specialty_id")).Preload(clause.Associations).Find(&entities)
	utilities.Show(c, http.StatusOK, "success", entities)
}

func GetDoctorsBySpecialty(c *gin.Context) {
	var payload struct {
		SearchParam string `json:"search_param" binding:"required"`
		SearchVal   string `json:"search_val" binding:"required"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	var entities []models.Doctor
	db := database.DB
	searchID := payload.SearchVal
	if payload.SearchParam == "slug" {
		searchID = appointments.GetEntityIDBySlug(models.Specialty{}, payload.SearchVal)
	}
	db.Where("specialty_id = ?", searchID).Preload(clause.Associations).Find(&entities)
	utilities.Show(c, http.StatusOK, "success", entities)
}

func GetDoctorDetails(c *gin.Context) {
	var payload struct {
		SearchParam string `json:"search_param" binding:"required"`
		SearchVal   string `json:"search_val" binding:"required"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	var entities []models.Doctor
	db := database.DB
	db.Where(payload.SearchParam+" = ?", payload.SearchVal).Preload(clause.Associations).First(&entities)
	utilities.Show(c, http.StatusOK, "success", entities)
}
