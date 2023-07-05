package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Quotation struct {
	gorm.Model
	types.CModel
	LeadDetailID uint
	ReferenceNo  string `gorm:"type:varchar(20)"`
	Date         string `gorm:"type:varchar(08)"`
	ProductType  string `gorm:"type:varchar(05)"`
	ProductCode  string `gorm:"type:varchar(05)"`
	Notes        string `gorm:"type:longtext"`
}
