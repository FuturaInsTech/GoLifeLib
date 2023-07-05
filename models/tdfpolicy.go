package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type TDFPolicy struct {
	gorm.Model
	types.CModel
	TDFType       string `gorm:"type:varchar(20)"`
	PolicyID      uint
	EffectiveDate string `gorm:"type:varchar(08);"`
	Seqno         uint16
}
