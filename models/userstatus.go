package models

import (
	"gorm.io/gorm"
)

type UserStatus struct {
	gorm.Model
	UserStatusShortName string `gorm:"type:varchar(8)"`
	UserStatusLongName  string `gorm:"type:varchar(80)"`
	Users               []User
}
