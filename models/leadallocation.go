package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type LeadAllocation struct {
	gorm.Model
	types.CModel

	Office          string `gorm:"type:varchar(05)"`
	SalesManager    string `gorm:"type:varchar(10)"`
	AgencyID        uint
	AllocationDate  string `gorm:"type:varchar(08)"`
	Priority        uint
	Quality         string `gorm:"type:varchar(01)"`
	AppointmentDate string `gorm:"type:varchar(08)"`
	LeadAllocStatus string `gorm:"type:varchar(05)"`
	ProductType     string `gorm:"type:varchar(05)"`
	ProductCode     string `gorm:"type:varchar(05)"`
	NoofAppointment uint
	ClosureStatus   string `gorm:"type:varchar(05)"`
	ClosureDate     string `gorm:"type:varchar(08)"`
	ExtractionDate  string `gorm:"type:varchar(08)"`
	LeadDetailID    uint
	ReceivedDate    string `gorm:"type:varchar(08)"`
	LeadChannelID   uint
	CampaignCode    string `gorm:"type:varchar(05)"`
	ClientID        uint
	ClientName      string `gorm:"type:varchar(50)"`
}
