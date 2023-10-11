package utilities

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var ErrrsList []string

func DatabaseErrorHandler(err error, tag string) string {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Sprintf("%s is not found. Check and try again", tag)
	} else if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint ") {
		return fmt.Sprintf("%s already exist. Check and try again", tag)
	} else {
		return err.Error()
	}
}
