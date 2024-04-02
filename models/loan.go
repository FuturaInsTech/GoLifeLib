package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Loan struct {
	gorm.Model
	types.CModel
	PolicyID        uint
	LoanSeqNumber   uint
	BenefitID       uint
	PProduct        string `gorm:"type:varchar(03)"`
	BCoverage       string `gorm:"type:varchar(05)"` //Q0006
	ClientID        uint
	TranDate        string `gorm:"type:varchar(08)"`
	TranNumber      uint
	LoanEffDate     string `gorm:"type:varchar(08)"`
	LoanType        string `gorm:"type:varchar(02)"`
	LoanStatus      string `gorm:"type:varchar(02)"`
	LoanCurrency    string `gorm:"type:varchar(03)"`
	LoanAmount      float64
	LoanIntRate     float64
	LoanIntType     string `gorm:"type:varchar(02)"`
	LastCapAmount   float64
	LastCapDate     string `gorm:"type:varchar(08)"`
	NextCapDate     string `gorm:"type:varchar(08)"`
	LastIntBillDate string `gorm:"type:varchar(08)"`
	NextIntBillDate string `gorm:"type:varchar(08)"`
	StampDuty       float64
	LoanBills       []LoanBill
}
