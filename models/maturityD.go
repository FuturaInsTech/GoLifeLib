package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type MaturityD struct {
	gorm.Model
	types.CModel
	MaturityHID         uint
	PolicyID            uint
	ClientID            uint
	BenefitID           uint
	BCoverage           string `gorm:"type:varchar(05)"` //Q0006
	BSumAssured         uint64
	MaturityAmount      float64
	RevBonus            float64
	AddlBonus           float64
	TerminalBonus       float64
	InterimBonus        float64
	LoyaltyBonus        float64
	OtherAmount         float64
	AccumDividend       float64
	AccumDivInt         float64
	TotalFundValue      float64
	TotalMaturityAmount float64
}
