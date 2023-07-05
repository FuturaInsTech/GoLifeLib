package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type GlType struct {
	gorm.Model
	types.CModel
	GlTypeCode        string `gorm:"type:varchar(01);primaryKey"`
	GlTypeDescription string `gorm:"type:varchar(10)"`
	AccountCodes      []AccountCode
}
