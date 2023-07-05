package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type GlMove struct {
	gorm.Model
	types.CModel
	GlRdocno       string `gorm:"type:varchar(20)"`
	GlRldgAcct     string `gorm:"type:varchar(30)"`
	GlCurry        string `gorm:"type:varchar(03)"`
	GlAmount       float64
	ContractCurry  string `gorm:"type:varchar(03)"`
	ContractAmount float64
	AccountCodeID  uint
	AccountCode    string `gorm:"type:varchar(30)"`
	GlSign         string `gorm:"type:varchar(01)"`
	SequenceNo     uint64
	CurrencyRate   float64
	CurrentDate    string `gorm:"type:varchar(08)"`
	EffectiveDate  string `gorm:"type:varchar(08)"`
	ReconciledDate string `gorm:"type:varchar(08)"`
	ExtractedDate  string `gorm:"type:varchar(30)"`
	HistoryCode    string `gorm:"type:varchar(05)"`
}
