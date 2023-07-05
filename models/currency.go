package models

import "gorm.io/gorm"

type Currency struct {
	gorm.Model
	CurrencyShortName string `gorm:"type:varchar(03)"`
	CurrencyLongName  string `gorm:"type:varchar(50)"`
	Companies         []Company
}
