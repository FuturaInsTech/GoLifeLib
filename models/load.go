package models

import "gorm.io/gorm"

type Load struct {
	gorm.Model
	StartDate  string `gorm:"type:varchar(8)"`
	EndDate    string `gorm:"type:varchar(8)"`
	Amount     float64
	CoverageID uint
}
