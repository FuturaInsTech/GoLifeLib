package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type PayingAuthority struct {
	gorm.Model
	types.CModel
	ClientID        uint
	PaName          string `gorm:"type:varchar(50)"`
	PaType          string `gorm:"type:varchar(01)"` //P0050
	StartDate       string `gorm:"type:varchar(08)"`
	EndDate         string `gorm:"type:varchar(08)"`
	PaStatus        string `gorm:"type:varchar(02)"` // P0050
	PaBillSummaries []PaBillSummary
}
