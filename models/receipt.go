package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

// ------13 Jan 2024
// ------RECEIPT TABLE IS REDESIGNED TO ACCEPT MULTIPLE TYPES OF RECEIPTS
//------ SEE BELOW THE NEW STRUCT
//------ NOTES: DROP EXISTING DB TABLE AND REGENERATE THE NEWLY DESIGNED TABLE COLUMNS

type Receipt struct {
	gorm.Model
	types.CModel
	Branch           string `gorm:"type:varchar(05)"` // P0017
	CurrentDate      string `gorm:"type:varchar(8)"`
	AccCurry         string `gorm:"type:varchar(03)"`
	AccAmount        float64
	ClientID         uint
	AddressID        uint
	DateOfCollection string `gorm:"type:varchar(08)"`
	ReceiptFor       string `gorm:"type:varchar(05)"` // P0050
	ReceiptRefNo     uint
	ReceiptCurry     string `gorm:"type:varchar(03)"`
	ReceiptAmount    float64
	ReceiptDueDate   string `gorm:"type:varchar(08)"`
	ReconciledDate   string `gorm:"type:varchar(08)"`
	BankIFSC         string `gorm:"type:varchar(50)"` // Client Bank Code
	BankAccountNo    string `gorm:"type:varchar(50)"` // Client Bank Account
	BankReferenceNo  string `gorm:"type:varchar(40)"`
	TypeOfReceipt    string `gorm:"type:varchar(05)"` // P0030
	InsurerBankIFSC  string `gorm:"type:varchar(50)"`
	InsurerBankAccNo string `gorm:"type:varchar(50)"`
}
