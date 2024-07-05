package clients

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

func AddClient(c *gin.Context) {
	var payload struct {
		FirstName   string `json:"firstname" binding:"required"`
		MiddleName  string `json:"middlename"`
		LastName    string `json:"lastname" binding:"required"`
		Email       string `json:"emailaddress"`
		Phone       string `json:"phone" binding:"required"`
		Language    string `json:"language" binding:"required"`
		Gender      string `json:"gender" binding:"required"`
		Address     string `json:"address" binding:"required"`
		Country     string `json:"country" binding:"required"`
		DateOfBirth string `json:"dob" binding:"required"`
		State       string `json:"county"`
		City        string `json:"subcounty"`
		PostalCode  string `json:"postal"`
		ZipCode     string `json:"zip"`
		People      []struct {
			FullName     string `json:"fullname"`
			Phone        string `json:"phoneno"`
			Relationship string `json:"relation"`
		} `json:"nextofkin"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	//database connection
	db := database.DB

	var u models.Client
	result := db.Where("email = ?", payload.Email).First(&u)
	if result.RowsAffected > 0 {
		utilities.ShowError(c, http.StatusOK, append(utilities.ErrrsList, "Client with this email already exist"))
		return
	}
	client := models.Client{
		FirstName:   payload.FirstName,
		MiddleName:  &payload.MiddleName,
		LastName:    payload.LastName,
		Email:       &payload.Email,
		Phone:       payload.Phone,
		GenderID:    payload.Gender,
		LanguageID:  payload.Language,
		DateOfBirth: payload.DateOfBirth,
		Address:     payload.Address,
		CountryID:   payload.Country,
		State:       &payload.State,
		City:        &payload.City,
		PostalCode:  &payload.PostalCode,
		ZipCode:     &payload.ZipCode,
		CreatedBy:   middlewares.GetAuthUserID(c),
	}
	if err := db.Create(&client).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	for _, people := range payload.People {
		people := models.ClientPeople{
			ClientID:     client.ID,
			FullName:     people.FullName,
			Phone:        people.Phone,
			Relationship: people.Relationship,
			CreatedBy:    middlewares.GetAuthUserID(c),
		}
		if err := db.Create(&people).Error; err != nil {
			utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	utilities.Show(c, http.StatusOK, "Client details added successfully", payload)
}

func ListClients(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.ClientsFilter)
	db := database.DB
	var entities []models.Client
	if (filters.ClientsFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "clients"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Joins("LEFT JOIN genders ON genders.id=clients.gender_id").
			Joins("LEFT JOIN languages ON languages.id=clients.language_id").
			Joins("LEFT JOIN countries ON countries.id=clients.country_id").
			Where("clients.first_name ILIKE ? OR clients.middle_name ILIKE ? OR clients.last_name ILIKE ? OR clients.email ILIKE ? OR clients.phone ILIKE ? OR genders.title ILIKE ? OR languages.title ILIKE ? OR countries.country_name ILIKE ?", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "clients"))
			return
		}
	}
	services.PaginationResponse(db, c, http.StatusOK, "clients", entities, models.Client{})

}
