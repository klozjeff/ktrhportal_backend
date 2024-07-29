package dashboard

import (
	"ktrhportal/database"
	"ktrhportal/models"
	"ktrhportal/utilities"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Overview(c *gin.Context) {
	db := database.DB
	clientCounts := clients(db)
	appointmentCounts := appointments(db)
	providers := providers(db)
	services := services(db)
	utilities.Show(c, http.StatusOK, "overview", map[string]interface{}{
		"clients":      clientCounts,
		"appointments": appointmentCounts,
		"providers":    providers,
		"services":     services,
	})
}

func clients(db *gorm.DB) map[string]int64 {
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	sixtyDaysAgo := now.AddDate(0, 0, -60)

	var last30DaysCount, previous30DaysCount int64

	db.Model(&models.Client{}).Where("created_at BETWEEN ? AND ?", thirtyDaysAgo, now).Count(&last30DaysCount)
	db.Model(&models.Client{}).Where("created_at BETWEEN ? AND ?", sixtyDaysAgo, thirtyDaysAgo).Count(&previous30DaysCount)

	return map[string]int64{
		"last30Days":     last30DaysCount,
		"previous30Days": previous30DaysCount,
	}
}

func appointments(db *gorm.DB) map[string]int64 {
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	sixtyDaysAgo := now.AddDate(0, 0, -60)

	var last30DaysCount, previous30DaysCount int64

	db.Model(&models.Appointment{}).Where("created_at BETWEEN ? AND ?", thirtyDaysAgo, now).Count(&last30DaysCount)
	db.Model(&models.Appointment{}).Where("created_at BETWEEN ? AND ?", sixtyDaysAgo, thirtyDaysAgo).Count(&previous30DaysCount)

	return map[string]int64{
		"last30Days":     last30DaysCount,
		"previous30Days": previous30DaysCount,
	}
}

func providers(db *gorm.DB) int64 {
	var count int64
	db.Model(&models.Provider{}).Count(&count)
	return count
}

func services(db *gorm.DB) int64 {
	var count int64
	db.Model(&models.Service{}).Count(&count)
	return count
}
