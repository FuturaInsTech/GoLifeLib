package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type PolBill struct {
	gorm.Model
	types.CModel
	PolicyID       uint
	ClientID       uint
	BankID         uint
	PBillCurr      string `gorm:"type:varchar(03)"` //P0023
	PolStatus      string `gorm:"type:varchar(05)"` //P0024 READ ONLY
	BillType       string `gorm:"type:varchar(05)"` //P0055
	BankCode       string `gorm:"type:varchar(50)"` //IFSC
	BankAccountNo  string `gorm:"type:varchar(50)"` // Bank Account No of Payer/Agent/Claiment
	InstalmentPrem float64
	CreationDate   string `gorm:"type:varchar(08)"`
	ExtractionDate string `gorm:"type:varchar(08)"`
	BillDate       string `gorm:"type:varchar(08)"`
	PaidToDate     string `gorm:"type:varchar(08)"`
	BillFreq       string `gorm:"type:varchar(08)"`
	PayeeName      string `gorm:"type:varchar(90)"`
}
