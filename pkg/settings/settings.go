package settings

import (
	"ktrhportal/database"
	"ktrhportal/middlewares"
	"ktrhportal/models"
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
