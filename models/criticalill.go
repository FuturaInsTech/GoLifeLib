package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type CriticalIllness struct {
	gorm.Model
	types.CModel
	EffectiveDate string `gorm:"type:varchar(8)"`
	IncidentDate  string `gorm:"type:varchar(8)"`
	PolicyID      uint
	BenefitID     uint
	CriticalType  string `gorm:"type:varchar(5)"`
	BSumAssured   uint64
	BStatusCode   string `gorm:"type:varchar(2)"`
	ApprovalFlag  string `gorm:"type:varchar(2)"`
}
