package auth

import (
	"ktrhportal/database"
	"ktrhportal/middlewares"
	"ktrhportal/models"
	"ktrhportal/utilities"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func Register(c *gin.Context) {
	var payload struct {
		Name     string `form:"name" json:"name" binding:"required"`
		Email    string `form:"email" json:"email" binding:"required"`
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
		Phone    string `form:"phone" json:"phone"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		var errorlist []string
		errorlist = append(errorlist, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, errorlist)
		return
	}
	role := GetEntityIDBySlug(models.Role{}, "admin")
	user := models.User{
		Name:            payload.Name,
		Email:           payload.Email,
		Phone:           payload.Phone,
		Username:        payload.Username,
		Password:        utilities.GeneratePassword(payload.Password),
		AccountStatusID: GetEntityIDBySlug(models.AccountStatus{}, "active"),
		RoleID:          role,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	//update last login
	database.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("last_logged_in", time.Now())
	token, _ := middlewares.GenerateToken(user.ID.String(), role)
	refreshToken, _ := middlewares.GenerateRefreshToken(user.ID.String(), role)

	var newUser models.User
	database.DB.Where("id=?", user.ID).
		Preload(clause.Associations).
		First(&newUser)
	utilities.Show(c, http.StatusOK, "User account created successfully", map[string]interface{}{
		"token":         token,
		"refresh_token": refreshToken,
		"username":      newUser.Username,
		"role":          newUser.Role.Title,
	})
}

func Login(c *gin.Context) {
	var payload struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}
	if validationError := c.ShouldBindJSON(&payload); validationError != nil {
		utilities.ErrrsList = append(utilities.ErrrsList, utilities.Validate(validationError)...)
		utilities.ShowError(c, http.StatusBadRequest, utilities.ErrrsList)
		return
	}
	var user models.User
	db := database.DB
	if err := db.Where("username = ?", payload.Username).Preload(clause.Associations).First(&user).Error; err != nil {
		utilities.ShowError(c, http.StatusOK, append(utilities.ErrrsList, "Wrong username or password"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		utilities.ShowError(c, http.StatusOK, append(utilities.ErrrsList, "Wrong username or password"))
		return
	}
	token, _ := middlewares.GenerateToken(user.ID.String(), user.RoleID)
	refreshToken, _ := middlewares.GenerateRefreshToken(user.ID.String(), user.RoleID)

	utilities.SetCookie(c, "token", token, time.Now().Add(time.Hour*1))

	utilities.Show(c, http.StatusOK, "Logged in successfully", map[string]interface{}{
		"token":         token,
		"refresh_token": refreshToken,
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"phone":         user.Phone,
		"username":      user.Username,
		"role":          user.Role.Title,
		"status":        user.AccountStatus.Slug,
	})

	/*c.JSON(http.StatusOK, gin.H{
		"token": token,
	})*/
}

func GetEntityIDBySlug(model interface{}, slug string) string {
	var id string
	database.DB.Model(&model).Where("slug=?", slug).Select("id").First(&id)
	return id
}

func CurrentUser(c *gin.Context) {
	var user models.User
	db := database.DB
	if err := db.Where("id = ?", middlewares.GetAuthUserID(c)).Preload(clause.Associations).First(&user).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "me", user)
}
