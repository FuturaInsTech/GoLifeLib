package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type GlType struct {
	gorm.Model
	types.CModel
	GlTypeCode        string `gorm:"type:varchar(01);primaryKey"`
	GlTypeDescription string `gorm:"type:varchar(10)"`
	AccountCodes      []AccountCode
}
