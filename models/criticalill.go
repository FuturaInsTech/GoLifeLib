package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type CriticalIllness struct {
	gorm.Model
	types.CModel

	PolicyID      uint
	BenefitID     uint
	CriticalType  string `gorm:"type:varchar(10)"`
	BSumAssured   uint64
	EffectiveDate string `gorm:"type:varchar(8)"`
	IncidentDate  string `gorm:"type:varchar(8)"`
	ReceivedDate  string `gorm:"type:varchar(8)"`
	PaidToDate    string `gorm:"type:varchar(8)"`
	BStatusCode   string `gorm:"type:varchar(2)"`
	ApprovalFlag  string `gorm:"type:varchar(2)"`
	ClaimAmount   uint64
	Percentage    float64
}
