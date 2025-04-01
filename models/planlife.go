package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type PlanLife struct {
	gorm.Model
	types.CModel
	PolicyID         uint
	BenefitID         uint
	BenefitPlan 	string `gorm:"type:varchar(10)"`
	ClientID         uint
	ClientRelcode 	string `gorm:"type:varchar(05)"`
	ClientReldesc	string `gorm:"type:varchar(20)"`
	BenefitPlanSA	float64
}