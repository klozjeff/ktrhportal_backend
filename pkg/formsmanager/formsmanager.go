package formsmanager

import (
	"encoding/json"
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

func AddTemplate(c *gin.Context) {
	var payload models.FormTemplate
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	db := database.DB
	if err := db.Create(&payload).Error; err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to create template: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Form template added successfully", payload)

}

func ListTemplates(c *gin.Context) {
	// Retrieve query parameters
	input := c.Request.Context().Value(httpin.Input).(*filters.FormsManagerFilter)
	db := database.DB
	var entities []models.FormTemplate
	if (filters.FormsManagerFilter{}) == *input {
		if err := db.
			Scopes(services.Paginate(c)).
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "form_templates"))
			return
		}
	} else if input.Global != "" {
		if err := db.
			Scopes(services.Paginate(c)).
			Where("form_templates.name ILIKE ?", "%"+input.Global+"%").
			Preload(clause.Associations).
			Find(&entities).Error; err != nil {
			utilities.ShowMessage(c, http.StatusOK, utilities.DatabaseErrorHandler(err, "form_templates"))
			return
		}
	}
	services.PaginationResponse(db, c, http.StatusOK, "form_templates", entities, models.Service{})
}

func GetTemplateDetails(c *gin.Context) {
	db := database.DB
	var template models.FormTemplate
	if err := db.
		Where("id = ?", c.Param("id")).
		Preload(clause.Associations).
		First(&template).Error; err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}
	utilities.Show(c, http.StatusOK, "success", template)
}

func UpdateTemplate(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	db := database.DB
	var template models.FormTemplate

	// First, find the record
	if err := db.Where("id = ?", payload["id"]).First(&template).Error; err != nil {
		utilities.ShowMessage(c, http.StatusNotFound, "Template not found")
		return
	}

	// Then, update it
	if err := db.Model(&template).Updates(payload).Error; err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	utilities.ShowMessage(c, http.StatusOK, "Form template updated successfully")
}

func FormSubmmision(c *gin.Context) {
	var payload models.FormSubmission
	if err := c.ShouldBindJSON(&payload); err != nil {
		utilities.ShowMessage(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	db := database.DB
	if err := db.Create(&payload).Error; err != nil {
		utilities.ShowMessage(c, http.StatusInternalServerError, "Failed to submit form data: "+err.Error())
		return
	}

	utilities.Show(c, http.StatusOK, "Form details submitted successfully", payload)

}

func ListFormSubmissions(c *gin.Context) {
	db := database.DB
	var submissions []models.FormSubmission
	result := db.Preload(clause.Associations).Find(&submissions).Find(&submissions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching submissions"})
		return
	}

	for i, submission := range submissions {
		var submissionData map[string]interface{}
		json.Unmarshal([]byte(submission.Data), &submissionData)

		var templateFields []map[string]interface{}
		json.Unmarshal([]byte(submission.Template.Fields), &templateFields)

		fieldMap := make(map[string]string)
		for _, field := range templateFields {
			id, _ := field["id"].(string)
			label, _ := field["label"].(string)
			fieldMap[id] = label
		}

		labeledData := make(map[string]interface{})
		for id, value := range submissionData {
			if label, ok := fieldMap[id]; ok {
				labeledData[label] = value
			} else {
				labeledData[id] = value // Fallback to ID if label not found
			}
		}

		// Convert labeled data back to JSON string
		labeledDataJSON, _ := json.Marshal(labeledData)
		submissions[i].Data = string(labeledDataJSON)
	}

	//c.JSON(http.StatusOK, submissions)
	utilities.Show(c, http.StatusOK, "submissions", submissions)
}
