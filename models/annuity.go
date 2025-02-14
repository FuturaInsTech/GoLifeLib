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
	BenefitID    uint
	BCoverage    string  `gorm:"type:varchar(05)"` //Q0006
	AnnuityM     string  `gorm:"type:varchar(05)"` //P0050
	AnnSA        float64 // Purpose Amount or Contribution
	CntCurr      string  `gorm:"type:varchar(03)"` //INR
	PayCurr      string  `gorm:"type:varchar(03)"` //USD
	AnnFreq      string  `gorm:"type:varchar(02)"` //12
	AnnEndDate   string  `gorm:"type:varchar(08)"` //99999999
	AnnStartDate string  `gorm:"type:varchar(08)"` //20230101
	AnnCurrDate  string  `gorm:"type:varchar(08)"` //20250201
	AnnNxtDate   string  `gorm:"type:varchar(08)"` //20250301
	AnnPecentage float64
	AnnAmount    float64 // Annuity Amount
	Paymentno    uint    // Payment Table ID
	Paystatus    string  `gorm:"type:varchar(02)"` // PR - Processed , PN - Pending
}
