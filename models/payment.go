package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	types.CModel
	Branch           string `gorm:"type:varchar(05)"` // P0017
	CurrentDate      string `gorm:"type:varchar(8)"`
	AccCurry         string `gorm:"type:varchar(03)"`
	AccAmount        float64
	PolicyID         uint
	ClientID         uint
	AddressID        uint
	DateOfPayment    string `gorm:"type:varchar(08)"`
	PaymentAccount   string `gorm:"type:varchar(50)"` // P0050
	ReconciledDate   string `gorm:"type:varchar(08)"`
	BankAccountNo    string `gorm:"type:varchar(50)"` // Client Bank Account
	BankReferenceNo  string `gorm:"type:varchar(40)"`
	TypeOfPayment    string `gorm:"type:varchar(05)"` // P0030
	BankIFSC         string `gorm:"type:varchar(50)"` // Client Bank Code
	InsurerBankIFSC  string `gorm:"type:varchar(50)"`
	InsurerBankAccNo string `gorm:"type:varchar(50)"`
	Status           string `gorm:"type:varchar(02)"`
	MakerUserID      uint
	CheckerUserID    uint
	Reason           string `gorm:"type:varchar(100)"`
}
