package models

import "gorm.io/gorm"

type Coverage struct {
	gorm.Model
	CoverageType         string `gorm:"type:varchar(10)"`
	StartDate            string `gorm:"type:varchar(8)"`
	EndDate              string `gorm:"type:varchar(8)"`
	Term                 string `gorm:"type:varchar(3)"`
	Status               string `gorm:"type:varchar(10)"`
	LifeId               string `gorm:"type:varchar(20)"`
	ContHeaderContractId string `gorm:"type:varchar(20);"`
	Loads                []Load
}
