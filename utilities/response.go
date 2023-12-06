package utilities

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Show(c *gin.Context, code int, msg interface{}, data interface{}) {
	c.JSON(code, gin.H{
		"success":     true,
		"status_code": code,
		"message":     msg,
		"data":        data,
	})
}
func ShowMessage(c *gin.Context, code int, msg interface{}) {
	c.JSON(code, gin.H{
		"success":     code == http.StatusOK,
		"status_code": code,
		"message":     msg,
	})
}
func ShowError(c *gin.Context, code int, errors []string) {
	c.JSON(code, gin.H{
		"success":     false,
		"status_code": code,
		"error":       errors,
	})
}
