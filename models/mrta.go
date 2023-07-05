package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Mrta struct {
	gorm.Model
	types.CModel
	PolicyID       uint
	ClientID       uint
	Pproduct       string `gorm:"type:varchar(05)"`
	BStartDate     string `gorm:"type:varchar(08)"`
	BCoverage      string `gorm:"type:varchar(05)"`
	BenefitID      uint
	BTerm          uint
	PremPayingTerm uint
	BSumAssured    float64
	Interest       float64
	InterimPeriod  float64
}
