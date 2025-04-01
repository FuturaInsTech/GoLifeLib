package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type PlanLife struct {
	gorm.Model
	types.CModel
	PolicyID       uint
	BenefitID      uint
	BenefitPlan    string `gorm:"type:varchar(10)"`
	ClientID       uint
	ClientRelcode  string `gorm:"type:varchar(05)"`
	ClientReldesc  string `gorm:"type:varchar(20)"`
	PGender        string `gorm:"type:varchar(01)"`
	PDOB           string `gorm:"type:varchar(08)"`
	PMortality     string `gorm:"type:varchar(01)"`
	PStatus        string `gorm:"type:varchar(02)"`
	PAge           uint
	PSmoker        string `gorm:"type:varchar(1)"`
	PStartDate     string `gorm:"type:varchar(08)"`
	PSumAssured    uint64
	PBasAnnualPrem float64    // Annualized Premium Before Applying Discount and Factor
	PLoadPrem      float64    // Loaded Premium
	PDiscPrem      float64    // Total Premium Discount
	PPrem          float64    // Total Premium
	PDiscountType  [5]string  `gorm:"type:varchar(01)"` // Discount Type
	PDiscountPrem  [5]float64 // Premium Discount
}
