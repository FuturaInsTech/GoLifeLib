package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type CampaignComp struct {
	gorm.Model
	types.CModel
	CampaignID   uint
	CampaignCode string `gorm:"type:varchar(05)"`
	Fee          string `gorm:"type:varchar(05)"`
	Basis        string `gorm:"type:varchar(08)"`
	MinLead      uint
	StartDate    string `gorm:"type:varchar(08)"`
	EndDate      string `gorm:"type:varchar(08)"`
}
