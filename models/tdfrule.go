package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
)

type TDFRule struct {
	types.CModel
	CompanyID uint
	Seqno     uint16 `gorm:"primaryKey;"`
	TDFType   string `gorm:"type:varchar(20)"`
}
