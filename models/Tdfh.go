package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Tdfh struct {
	gorm.Model
	types.CModel
	PolicyID      uint
	EffectiveDate string `gorm:"type:varchar(08);"`
}
