package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Benefit struct {
	gorm.Model
	types.CModel
	ClientID         uint
	PolicyID         uint
	BStartDate       string `gorm:"type:varchar(08)"`
	BRiskCessDate    string `gorm:"type:varchar(08)"`
	BPremCessDate    string `gorm:"type:varchar(08)"`
	BTerm            uint
	BPTerm           uint
	BRiskCessAge     uint
	BPremCessAge     uint
	BBasAnnualPrem   float64 // Annualized Premium Before Applying Discount and Factor
	BLoadPrem        float64 // Loaded Premium
	BCoverage        string  `gorm:"type:varchar(05)"` //Q0006
	BSumAssured      uint64
	BPrem            float64 // Total Premium
	BGender          string  `gorm:"type:varchar(01)"`
	BDOB             string  `gorm:"type:varchar(08)"`
	BMortality       string  `gorm:"type:varchar(01)"`
	BStatus          string  `gorm:"type:varchar(02)"`
	BAge             uint
	BRerate          string `gorm:"type:varchar(08)"`
	BonusDate        string `gorm:"type:varchar(08)"`
	IlpMortality     float64
	IlpMortalityDate string `gorm:"type:varchar(08)"`
	IlpFee           float64
	IlpFeeDate       string `gorm:"type:varchar(08)"`
	Extras           []Extra
	SurvBs           []SurvB
	MedReqs          []MedReq
	DeathDs          []DeathD
	SaChanges        []SaChange
	Mrtas            []Mrta
	CriticalIllnesss []CriticalIllness
	IBenefits        []IBenefit
	IlpFunds         []IlpFund
	IlpTransactions  []IlpTransaction
	IlpSummaries     []IlpSummary
	IlpFundSwitches  []IlpFundSwitch
}
