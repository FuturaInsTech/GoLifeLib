package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Addcomponent struct {
	gorm.Model
	types.CModel
	Select        string `gorm:"type:varchar(01)"`
	PolicyID      uint
	ClientID      uint
	BCoverage     string `gorm:"type:varchar(05)"` //Q0006
	BStartDate    string `gorm:"type:varchar(08)"`
	BSumAssured   uint64
	BTerm         uint
	BPTerm        uint
	BPrem         float64 // Total Premium
	BAnnualPrem   float64
	BGender       string `gorm:"type:varchar(01)"`
	BDOB          string `gorm:"type:varchar(08)"`
	Method        string `gorm:"type:varchar(10)"` // indicate whether it is a sa change or component add
	Frequency     string `gorm:"type:varchar(01)"`
	BAge          uint
	BRiskCessDate string `gorm:"type:varchar(08)"` // When rider date > basic, then it will be basic date
	BPremCessDate string `gorm:"type:varchar(08)"` // when rider date > basic, then it will be basic date
}

type Cola struct {
	gorm.Model
	types.CModel
	PolicyID         uint
	BenefitID        uint
	BCoverage        string `gorm:"type:varchar(05)"`
	ColaMethod       string `gorm:"type:varchar(06)"`
	BStartDate       string `gorm:"type:varchar(08)"` //20250525
	BEndDate         string `gorm:"type:varchar(08)"` //20450525
	BOriginalSA      uint64
	BOrginalPrem     float64
	BNewStartDate    string `gorm:"type:varchar(08)"`
	BNewSA           uint64
	BNewAnnualPrem   float64
	BNewModelPrem    float64
	BNewAge          uint
	Active           string `gorm:"type:varchar(01)"`
	Processed        string `gorm:"type:varchar(01)"`
	NewOrExist       string `gorm:"type:varchar(01)"`
	BNewTerm         uint
	BNewPTerm        uint
	CommissionAmount float64
}
