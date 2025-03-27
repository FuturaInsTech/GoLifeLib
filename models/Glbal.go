package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type GlBal struct {
	gorm.Model
	types.CModel
	GlRdocno       string `gorm:"type:varchar(20)"`
	GlRldgAcct     string `gorm:"type:varchar(30)"`
	GlAccountno    string `gorm:"type:varchar(30)"`
	ContractCurry  string `gorm:"type:varchar(03)"`
	ContractAmount float64
}

type PayOsBal struct {
	gorm.Model
	types.CModel
	GlRdocno       string `gorm:"type:varchar(20)"`
	GlRldgAcct     string `gorm:"type:varchar(30)"`
	GlAccountno    string `gorm:"type:varchar(30)"`
	ContractCurry  string `gorm:"type:varchar(03)"`
	PaymentNo      uint
	ContractAmount float64
}
