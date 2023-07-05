package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type LeadFollowup struct {
	gorm.Model
	types.CModel
	LeadDetailID      uint
	SeqNo             uint
	AppointmentDate   string `gorm:"type:varchar(08)"`
	AppointmentFlag   string `gorm:"type:varchar(02)"` //P0005
	ActualMeetingDate string `gorm:"type:varchar(08)"`
	ActionNote        string `gorm:"type:longtext"`
	ProgressStatus    string `gorm:"type:varchar(02)"` //P0004
	NextFollowupDate  string `gorm:"type:varchar(08)"`
	CountryCode       string `gorm:"type:varchar(03)"` // P0037
	PreferredDay      string `gorm:"type:varchar(01)"` // P0038
	PreferredTime     string `gorm:"type:varchar(01)"` // P0039
	AgencyID          uint
	SalesManager      string `gorm:"type:varchar(10)"`
	AllocationDate    string `gorm:"type:varchar(08)"`
	LeadAllocStatus   string `gorm:"type:varchar(05)"`
	ClosureStatus     string `gorm:"type:varchar(05)"`
	ClosureDate       string `gorm:"type:varchar(08)"`
	ClientID          uint
	ClientName        string `gorm:"type:varchar(50)"`
}
