package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type MaturityH struct {
	gorm.Model
	types.CModel
	PolicyID             uint
	ClientID             uint
	EffectiveDate        string `gorm:"type:varchar(08)"`
	MaturityDate         string `gorm:"type:varchar(08)"`
	Status               string `gorm:"type:varchar(02)"`
	BillDate             string `gorm:"type:varchar(08)"`
	PaidToDate           string `gorm:"type:varchar(08)"`
	Product              string `gorm:"type:varchar(05)"`
	AplAmount            float64
	LoanAmount           float64
	PolicyDepost         float64
	CashDeposit          float64
	RefundPrem           float64
	PremTolerance        float64
	TotalMaturityPayable float64
	AdjustedAmount       float64
	MaturityDs           []MaturityD
}
