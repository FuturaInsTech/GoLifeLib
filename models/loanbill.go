package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type LoanBill struct {
	gorm.Model
	types.CModel
	PolicyID        uint
	LoanSeqNumber   uint
	BenefitID       uint
	ClientID        uint
	LoanID          uint // Foreign Key of Loan
	TranNumber      uint
	PolStatus       string `gorm:"type:varchar(05)"` //P0024
	PaidToDate      string `gorm:"type:varchar(08)"`
	LoanBillCurr    string `gorm:"type:varchar(03)"`
	LoanType        string `gorm:"type:varchar(02)"`
	LoanBillDueDate string `gorm:"type:varchar(08)"`
	LoanIntAmount   float64
	PayerID         string `gorm:"type:varchar(30)"`
	ReceiptNo       uint
	ReceiptDate     string `gorm:"type:varchar(08)"`
	ReceiptAmount   float64
	CreditBank      string `gorm:"type:varchar(30)"`
	CreationDate    string `gorm:"type:varchar(30)"`
	ExtractionDate  string `gorm:"type:varchar(08)"`
	ReconciledFlg   string `gorm:"type:varchar(01)"`
	ReconciledDate  string `gorm:"type:varchar(08)"`
}
