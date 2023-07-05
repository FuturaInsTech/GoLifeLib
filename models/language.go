package models

import (
	"gorm.io/gorm"
)

type Language struct {
	gorm.Model
	LangShortName string `gorm:"type:varchar(2)"`
	LangLongName  string `gorm:"type:varchar(100)"`
	Users         []User

	Errors []Error
}
