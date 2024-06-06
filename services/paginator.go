package services

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int
	Page       int
	Sort       string
	TotalRows  int64
	TotalPages int
	Previous   int
	Next       int
}

var Limit int
var Page int
var Sort string

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		sortColumn := "created_at"
		sortTerm := "desc"
		limit := 30
		page := 1
		Page = 1
		Limit = 30
		query := c.Request.URL.Query()

		for key, value := range query {
			queryValue := value[len(value)-1]
			switch key {
			case "limit":
				limit, _ = strconv.Atoi(queryValue)
				Limit = limit
			case "page":
				page, _ = strconv.Atoi(queryValue)
				Page = page
			case "sortColumn":
				sortColumn = queryValue
			case "sortTerm":
				sortTerm = queryValue

			}

		}
		sort := fmt.Sprintf("%s %s", sortColumn, sortTerm)
		Sort = sort

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit).Order(sort)
	}
}

func PaginationResult(db *gorm.DB, c *gin.Context, data interface{}, model interface{}) Pagination {
	sortColumn := "created_at"
	sortTerm := "asc"
	limit := 30
	page := 1
	Page = 1
	Limit = 30
	query := c.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			Limit = limit
		case "page":
			page, _ = strconv.Atoi(queryValue)
			Page = page
		case "sortColumn":
			sortColumn = queryValue
		case "sortTerm":
			sortTerm = queryValue

		}

	}
	sort := fmt.Sprintf("%s %s", sortColumn, sortTerm)
	Sort = sort

	var pagination Pagination
	var totalRows int64
	db.Model(model).Count(&totalRows)
	pagination.TotalRows = totalRows
	pagination.Limit = limit
	pagination.Page = page
	pagination.Sort = sort
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
	pagination.Previous = GetPreviousPage(pagination.Page)
	pagination.Next = GetNextPage(pagination.Page, pagination.TotalPages)
	return pagination
}

func GetPreviousPage(page int) int {
	if page > 1 {
		return page - 1
	}
	return 0
}

func GetNextPage(page int, totalPages int) int {
	if page < totalPages {
		return page + 1
	}
	return 0
}

func PaginationResponse(db *gorm.DB, c *gin.Context, status int, message string, data interface{}, model interface{}) {
	var pagination Pagination = PaginationResult(db, c, data, model)
	c.JSON(http.StatusOK, gin.H{
		"message":      message,
		"status_code":  status,
		"success":      true,
		"data":         data,
		"sort":         pagination.Sort,
		"current_page": pagination.Page,
		"limit":        pagination.Limit,
		"previous":     pagination.Previous,
		"next":         pagination.Next,
		"total":        pagination.TotalRows,
		"pages":        pagination.TotalPages,
	})
}
