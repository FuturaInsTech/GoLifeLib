package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Annuity struct {
	gorm.Model
	types.CModel
	PolicyID     uint
	ClientID     uint
	AnnAmount    float64 // Annuity Amount
	AnnPecentage float64
	BenefitID    uint
	AnnType      string  `gorm:"type:varchar(05)"` //01    /P0050
	BCoverage    string  `gorm:"type:varchar(05)"` //Q0006
	AnnSA        float64 // Purpose Amount or Contribution
	AnnEndDate   string  `gorm:"type:varchar(08)"` //99999999
	AnnStartDate string  `gorm:"type:varchar(08)"` //20230101
	AnnCurrDate  string  `gorm:"type:varchar(08)"` //20250201
	AnnNxtDate   string  `gorm:"type:varchar(08)"` //20250301
	AnnFreq      string  `gorm:"type:varchar(02)"` //12
	CntCurr      string  `gorm:"type:varchar(03)"` //INR
	PayCurr      string  `gorm:"type:varchar(03)"` //USD
}
