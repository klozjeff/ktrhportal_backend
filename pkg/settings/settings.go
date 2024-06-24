package settings

import (
	"ktrhportal/database"
	"ktrhportal/filters"
	"ktrhportal/middlewares"
	"ktrhportal/models"
	"ktrhportal/pkg/appointments"
	"ktrhportal/services"
	"ktrhportal/utilities"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
func GetAppointmentStatuses(c *gin.Context) {
	db := database.DB
	var entities []models.AppointmentStatus

	if err := db.Find(&entities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusFound, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "appointment statuses", entities)
}

func GetEncounterStatuses(c *gin.Context) {
	db := database.DB
	var entities []models.EncounterStatus

	if err := db.Find(&entities).Error; err != nil {
		utilities.ShowMessage(c, http.StatusFound, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "encounter statuses", entities)
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
func GetCountries(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.CountriesFilter)
	db := database.DB
	var entities []models.Country
	if (filters.CountriesFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "countries"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Where("countries.country_name ILIKE ?", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "countries"))
			return
		}
	}
	services.PaginationResponse(db, c, http.StatusOK, "countries", entities, models.Country{})
}

func GetCounties(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.CountiesFilter)
	db := database.DB
	var entities []models.County
	if (filters.CountiesFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "counties"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Where("counties.name ILIKE ?", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "counties"))
			return
		}
	}
	services.PaginationResponse(db, c, http.StatusOK, "counties", entities, models.County{})
}

func GetSubCounties(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.SubCountiesFilter)
	db := database.DB
	var entities []models.SubCounty
	if (filters.SubCountiesFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "sub-counties"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Where("sub_counties.name ILIKE ?", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "sub-counties"))
			return
		}
	}
	services.PaginationResponse(db, c, http.StatusOK, "sub-counties", entities, models.SubCounty{})
}

/*
	func GetCounties(c *gin.Context) {
		db := database.DB
		var entities []models.County

		if err := db.Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusFound, err.Error())
			return
		}
		utilities.Show(c, http.StatusOK, "counties", entities)
	}
*/
func GetSubcounties(c *gin.Context) {
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

func AllSpecialties(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.SpecialtiesFilter)
	db := database.DB
	var entities []models.Specialty
	if (filters.SpecialtiesFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "specialties"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Where("specialties.name ILIKE ?", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "specialties"))
			return
		}
	}
	services.PaginationResponse(db, c, http.StatusOK, "specialties", entities, models.Specialty{})
	//utilities.Show(c, http.StatusOK, "specialties", entities)
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

func AllDoctors(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.DoctorsFilter)
	db := database.DB
	var entities []models.Doctor
	var scopes []func(*gorm.DB) *gorm.DB
	if (filters.DoctorsFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "doctors"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Joins("LEFT JOIN specialties ON specialties.id=doctors.specialty_id").
			Where("doctors.first_name ILIKE ? OR doctors.last_name ILIKE ? OR doctors.email ILIKE ? OR doctors.phone ILIKE ? OR specialties.name ILIKE ?", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "doctors"))
			return
		}
	} else {

		if input.SpecialtyID != "" {
			scopes = append(scopes, input.SpecialtyFilter)
		}
		if err := db.
			Scopes(scopes...).
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "doctors"))
			return
		}
	}

	services.PaginationResponse(db, c, http.StatusOK, "doctors", entities, models.Doctor{})

	//utilities.Show(c, http.StatusOK, "doctors", entities)
}
