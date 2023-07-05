package models

import (
	"github.com/shijith.chand/go-jwt/models/quotation"
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type Agency struct {
	gorm.Model
	types.CModel
	AgencyChannelSt string `gorm:"type:varchar(05)"` // P0017
	Office          string `gorm:"type:varchar(05)"` // P0018
	AgencySt        string `gorm:"type:varchar(05)"` // P0019
	LicenseNo         string `gorm:"type:varchar(20)"`
	LicenseStartDate  string `gorm:"type:varchar(08)"`
	LicenseEndDate    string `gorm:"type:varchar(08)"`
	Startdate         string `gorm:"type:varchar(08)"`
	EndDate           string `gorm:"type:varchar(08)"`
	TerminationReason string `gorm:"type:longtext"`
	ClientID          uint
	Aadhar            string `gorm:"type:varchar(020)"`
	Pan               string `gorm:"type:varchar(20)"`
	// AddressID         uint
	LeadAllocations []LeadAllocation
	BankID          uint
	Communications  []Communication
	Policies        []Policy
	QHeaders        []quotation.QHeader
}
