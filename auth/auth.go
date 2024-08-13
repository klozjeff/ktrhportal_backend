package auth

import (
	"fmt"
	"ktrhportal/database"
	"ktrhportal/middlewares"
	"ktrhportal/models"
	"ktrhportal/utilities"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func UserAccount(c *gin.Context) {
	var user models.User
	db := database.DB
	if err := db.Where("id = ?", middlewares.GetAuthUserID(c)).Preload(clause.Associations).First(&user).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "me", user)
}

func Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.Split(tokenString, "Bearer ")[1]
	if err := middlewares.InvalidateToken(tokenString); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.ShowMessage(c, http.StatusOK, "Successfully logged out")
}

func UpdateUserProfile(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Failed to retrieve file")
		return
	}

	db := database.DB
	var user models.User
	if err := db.Where("id = ?", middlewares.GetAuthUserID(c)).First(&user).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error)
		return
	}

	// Generate a random UUID and use it as the new file name
	fileID := uuid.New().String()
	fileExt := filepath.Ext(file.Filename)
	filename := fileID + fileExt

	// Ensure the upload directory exists
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Failed to save file")
		return
	}

	// Construct the URL to access the uploaded file
	fileURL := fmt.Sprintf("%s/api/v1/%s", utilities.GoDotEnvVariable("APP_URL"), dst)

	user.Profile = fileURL

	if err := db.Save(&user).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error)
		return
	}

	utilities.Show(c, http.StatusOK, "profile photo updated successfully", map[string]interface{}{
		"file_url": fileURL,
	})
}

func RemoveUserProfile(c *gin.Context) {
	db := database.DB
	var user models.User
	if err := db.Where("id = ?", middlewares.GetAuthUserID(c)).First(&user).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error)
		return
	}

	if user.Profile != "" {
		// Extract filename from the URL
		filename := path.Base(user.Profile)
		// Remove any URL encoding if present
		filename = strings.ReplaceAll(filename, "%20", " ")
		filePath := filepath.Join("uploads", filename)
		err := os.Remove(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("File does not exist: %s", filePath)
			} else if os.IsPermission(err) {
				log.Printf("Permission denied: %s", filePath)
				utilities.ShowMessage(c, http.StatusInternalServerError, "Permission denied when removing avatar file")
				return
			} else {
				log.Printf("Error removing file: %s. Error: %v", filePath, err)
				utilities.ShowMessage(c, http.StatusInternalServerError, fmt.Sprintf("Failed to remove avatar file: %v", err))
				return
			}
		} else {
			log.Printf("Successfully removed file: %s", filePath)
		}
	}

	// Update the database to remove the avatar reference
	if err := db.Model(&user).Update("Profile", "").Error; err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to update user data")
		return
	}
	utilities.ShowMessage(c, http.StatusOK, "Profile photo removed successfully")
}

func ServeUploadedFile(c *gin.Context) {
	filename := c.Param("filename")
	// Construct the file path
	filePath := filepath.Join("uploads", filename)
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	// Serve the file
	c.File(filePath)
}
