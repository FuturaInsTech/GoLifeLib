package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type LeadDetail struct {
	gorm.Model
	types.CModel
	CompanyID       uint
	LeadChannelID   uint
	OfficeCode      string `gorm:"type:varchar(05)"`
	ProviderName    string `gorm:"type:varchar(50)"`
	ReceivedDate    string `gorm:"type:varchar(08)"`
	CampaignCode    string `gorm:"type:varchar(05)"`
	ProductType     string `gorm:"type:varchar(05)"` // Q0005
	ProductCode     string `gorm:"type:varchar(05)"`
	ClientID        uint
	ClientName      string `gorm:"type:varchar(50)"`
	LeadFollowups   []LeadFollowup
	Quotations      []Quotation
	LeadAllocations []LeadAllocation
}
