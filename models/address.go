package models

import (
	"github.com/FuturaInsTech/GoLifeLib/models/quotation"
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	types.CModel
	AddressType      string `gorm:"type:varchar(05)"` //P0022
	AddressLine1     string `gorm:"type:varchar(50)"`
	AddressLine2     string `gorm:"type:varchar(50)"`
	AddressLine3     string `gorm:"type:varchar(50)"`
	AddressLine4     string `gorm:"type:varchar(50)"`
	AddressLine5     string `gorm:"type:varchar(50)"`
	AddressPostCode  string `gorm:"type:varchar(10)"`
	AddressState     string `gorm:"type:varchar(50)"`
	AddressCountry   string `gorm:"type:varchar(50)"`
	AddressStartDate string `gorm:"type:varchar(8)"`
	AddressEndDate   string `gorm:"type:varchar(8)"`
	ClientID         uint
	QHeaders         []quotation.QHeader
	// Agencies         []Agency
	Policies     []Policy
	MedProviders []MedProvider
	Receipts     []Receipt
}
