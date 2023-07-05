package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type LeadChannel struct {
	gorm.Model
	types.CModel
	CompanyID    uint
	ChannelCode  string `gorm:"type:varchar(02)"`
	ChannelDesc  string `gorm:"type:varchar(50)"`
	StartDate    string `gorm:"type:varchar(08)"`
	EndDate      string `gorm:"type:varchar(08)"`
	LeadAllocSt  string `gorm:"type:varchar(05)"` //P0015
	StatusReason string `gorm:"type:varchar(08)"`
	LeadDetails  []LeadDetail
	LeadAllocations []LeadAllocation
}
