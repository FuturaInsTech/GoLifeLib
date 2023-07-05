package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model
	types.CModel
	ChannelCode   string `gorm:"type:varchar(05)"`
	SourceName    string `gorm:"type:varchar(40)"`
	Province      string `gorm:"type:varchar(05)"`
	Region        string `gorm:"type:varchar(05)"`
	Office        string `gorm:"type:varchar(05)"`
	StartDate     string `gorm:"type:varchar(08)"`
	EndDate       string `gorm:"type:varchar(08)"`
	Status        string `gorm:"type:varchar(05)"`
	CampaignComps []CampaignComp
}
