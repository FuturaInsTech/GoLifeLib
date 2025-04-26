package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type CbUser struct {
	gorm.Model
	types.CModel

	ClientMobCode string `gorm:"type:varchar(05)"`
	ClientMobile  string `gorm:"type:varchar(20)"`
	ClientDob     string `gorm:"type:varchar(08)"`
	CbEmail       string `gorm:"type:varchar(100)"`
	CbEnabled     string `gorm:"type:varchar(01)"`
	Latitude      string `gorm:"type:varchar(100)"`
	Longitude     string `gorm:"type:varchar(100)"`
}

type CbLog struct {
	gorm.Model
	types.CModel
	CbID        uint
	UserMessage string `gorm:"type:longtext"`
	BotResponse string `gorm:"type:longtext"`
}
