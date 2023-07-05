package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type TBenefit struct {
	gorm.Model
	types.CModel
	ClientID       uint
	PolicyID       uint
	BStartDate     string `gorm:"type:varchar(08)"`
	BRiskCessDate  string `gorm:"type:varchar(08)"`
	BPremCessDate  string `gorm:"type:varchar(08)"`
	BTerm          uint
	BPTerm         uint
	BRiskCessAge   uint
	BPremCessAge   uint
	BBasAnnualPrem float64 // Annualized Premium Before Applying Discount and Factor
	BLoadPrem      float64 // Loaded Premium
	BCoverage      string  `gorm:"type:varchar(05)"` //Q0006
	BSumAssured    uint64
	BPrem          float64 // Total Premium
	BGender        string  `gorm:"type:varchar(01)"`
	BDOB           string  `gorm:"type:varchar(08)"`
	BMortality     string  `gorm:"type:varchar(01)"`
	BStatus        string  `gorm:"type:varchar(02)"`
	BAge           uint
	BRerate        string `gorm:"type:varchar(08)"`
	Extras         []Extra
}
