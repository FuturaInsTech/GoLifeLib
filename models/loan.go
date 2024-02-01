package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Loan struct {
	gorm.Model
	types.CModel
	PolicyID        uint
	BenefitID       uint
	PProduct        string `gorm:"type:varchar(03)"`
	BCoverage       string `gorm:"type:varchar(05)"` //Q0006
	ClientID        uint
	LoanSeqNumber   uint
	TranDate        string `gorm:"type:varchar(08)"`
	TranNumber      uint
	LoanEffDate     string `gorm:"type:varchar(08)"`
	LoanType        string `gorm:"type:varchar(02)"`
	LoanStatus      string `gorm:"type:varchar(02)"`
	LoanCurrency    string `gorm:"type:varchar(03)"`
	LoanAmount      string `gorm:"type:varchar(15)"`
	LoanIntRate     float64
	LoanIntType     string `gorm:"type:varchar(02)"`
	LastCapAmount   float64
	LastCapDate     string `gorm:"type:varchar(08)"`
	NextCapDate     string `gorm:"type:varchar(08)"`
	LastIntBillDate string `gorm:"type:varchar(08)"`
	NextIntBillDate string `gorm:"type:varchar(08)"`
	LoanBills       []LoanBill
}
