package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type SurvB struct {
	gorm.Model
	types.CModel
	Sequence      int
	PolicyID      uint
	BenefitID     uint
	EffectiveDate string `gorm:"type:varchar(08)"` //Due Date
	PaidDate      string `gorm:"type:varchar(08)"` //Paid Date
	Amount        float64
	SBPercentage  float64 //Percentage
}
