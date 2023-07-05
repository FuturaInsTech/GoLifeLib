package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type SaChange struct {
	gorm.Model
	types.CModel
	Select      string `gorm:"type:varchar(01)"`
	PolicyID    uint
	ClientID    uint
	BenefitID   uint
	BCoverage   string `gorm:"type:varchar(05)"` //Q0006
	BStartDate  string `gorm:"type:varchar(08)"`
	BSumAssured uint64
	BTerm       uint
	BPTerm      uint
	BPrem       float64 // Total Premium
	BGender     string  `gorm:"type:varchar(01)"`
	BDOB        string  `gorm:"type:varchar(08)"`
	NSumAssured uint64
	NTerm       uint
	NPTerm      uint
	NPrem       float64 // Instalment Premium
	NAnnualPrem float64
	Method      string `gorm:"type:varchar(10)"` // indicate whether it is a sa change or component add
	Frequency   string `gorm:"type:varchar(01)"`
}
