package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type UinMaster struct {
	gorm.Model
	types.CModel

	PlanCodeName        string `gorm:"type:varchar(10)"`
	PlanName            string `gorm:"type:varchar(50)"`
	GsvFactor           string `gorm:"type:varchar(50)"`
	GsvCashValue        string `gorm:"type:varchar(50)"`
	SsvFactor           string `gorm:"type:varchar(50)"`
	ProductType         string `gorm:"type:varchar(50)"`
	FlcEligibility      string `gorm:"type:varchar(50)"`
	SurrenderChargeRate uint64
}
