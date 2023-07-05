package models

import (
	"github.com/shijith.chand/go-jwt/types"
)

type TDFRule struct {
	types.CModel
	CompanyID uint
	Seqno     uint16 `gorm:"primaryKey;"`
	TDFType   string `gorm:"type:varchar(20)"`
}
