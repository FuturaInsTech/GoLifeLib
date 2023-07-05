package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Receipt struct {
	gorm.Model
	types.CModel
	Branch            string `gorm:"type:varchar(05)"` // P0017
	CurrentDate       string `gorm:"type:varchar(8)"`
	AccCurry          string `gorm:"type:varchar(03)"`
	AccAmount         float64
	PolicyID          uint
	ClientID          uint
	DateOfCollection  string `gorm:"type:varchar(08)"`
	ReconciledDate    string `gorm:"type:varchar(08)"`
	BankAccountNo     string `gorm:"type:varchar(40)"`
	BankReferenceNo   string `gorm:"type:varchar(40)"`
	TypeOfReceipt     string `gorm:"type:varchar(05)"` // P0030
	InstalmentPremium float64
	PaidToDate        string `gorm:"type:varchar(08)"`
	AddressID         uint
}
