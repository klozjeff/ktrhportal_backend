package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID      `json:"id" gorm:"primary_key"`
	Name            string         `json:"name"`
	Email           string         `json:"email"`
	Username        string         `gorm:"size:255;not null;unique" json:"username"`
	Password        string         `gorm:"size:255;not null;" json:"-"`
	Phone           string         `json:"phone"`
	RoleID          string         `json:"role_id"`
	Role            *Role          `json:"role"`
	AccountStatusID string         `json:"account_status_id"`
	AccountStatus   *AccountStatus `json:"status"`
	LastLogin       *time.Time     `json:"last_login"`
	CreatedAt       time.Time      `json:"-"`
	UpdatedAt       time.Time      `json:"-"`
	DeletedAt       gorm.DeletedAt `json:"-"`
}

func (user *User) BeforeCreate(scope *gorm.DB) error {

	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	var statusID string
	scope.Model(AccountStatus{}).Select("id").Where("slug=?", "active").First(&statusID)
	user.ID = uuid
	user.AccountStatusID = statusID
	return err
}
