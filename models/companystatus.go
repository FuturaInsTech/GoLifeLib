package models

import "gorm.io/gorm"

type CompanyStatus struct {
	gorm.Model
	ClientStatusShortName string `gorm:"type:varchar(8)"`
	ClientStatusLongName  string `gorm:"type:varchar(50)"`
	Companies             []Company
}
