package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type IlpPrice struct {
	gorm.Model
	types.CModel
	FundCode       string `gorm:"type:varchar(5)"` //P0050
	FundType       string `gorm:"type:varchar(2)"` //ACcummulated/INinitial/BOus P0050
	FundDate       string `gorm:"type:varchar(8)"`
	FundEffDate    string `gorm:"type:varchar(8)"`
	FundCurr       string `gorm:"type:varchar(3)"`
	FundBidPrice   float64
	FundOfferPrice float64
	FundSeqno      uint
	ApprovalFlag   string `gorm:"type:varchar(2)"`
}

type IlpSummary struct {
	gorm.Model
	types.CModel
	PolicyID  uint
	BenefitID uint
	FundCode  string  `gorm:"type:varchar(5)"` //P0050
	FundType  string  `gorm:"type:varchar(2)"` //P0050
	FundUnits float64 `gorm:"type:decimal(15,5);"`
	FundCurr  string  `gorm:"type:varchar(3)"`
}

type IlpAnnSummary struct {
	gorm.Model
	types.CModel
	PolicyID      uint
	BenefitID     uint
	FundCode      string  `gorm:"type:varchar(5)"` //P0050
	FundType      string  `gorm:"type:varchar(2)"` //P0050
	FundUnits     float64 `gorm:"type:decimal(15,5);"`
	FundCurr      string  `gorm:"type:varchar(3)"`
	EffectiveDate string  `gorm:"type:varchar(8)"`
}

type IlpTransaction struct {
	gorm.Model
	types.CModel
	PolicyID            uint
	BenefitID           uint
	FundCode            string `gorm:"type:varchar(5)"` //P0050
	FundType            string `gorm:"type:varchar(2)"` //ACcummulated/INinitial/BOus P0050
	TransactionDate     string `gorm:"type:varchar(8)"`
	FundEffDate         string `gorm:"type:varchar(8)"`
	FundAmount          float64
	FundCurr            string  `gorm:"type:varchar(3)"`
	FundUnits           float64 `gorm:"type:decimal(15,5);"`
	FundPrice           float64 `gorm:"type:decimal(15,5);"`
	CurrentOrFuture     string  `gorm:"type:varchar(1)"` //P0050
	OriginalAmount      float64
	ContractCurry       string  `gorm:"type:varchar(3)"`
	HistoryCode         string  `gorm:"type:varchar(5)"`
	InvNonInvFlag       string  `gorm:"type:varchar(2)"`
	InvNonInvPercentage float64 `gorm:"type:decimal(15,5);"`
	AccountCode         string  `gorm:"type:varchar(30)"`
	CurrencyRate        float64
	MortalityIndicator  string `gorm:"type:varchar(1)"`
	SurrenderPercentage float64
	Seqno               uint
	UlProcessFlag       string `gorm:"type:varchar(01)"` //P0050
	UlpPriceDate        string `gorm:"type:varchar(08)"`
	AllocationCategory  string `gorm:"type:varchar(02)"` //P0050  Denote PR/MP/FE etc
	AdjustedDate        string `gorm:"type:varchar(08)"`
}

type IlpFund struct {
	gorm.Model
	types.CModel
	PolicyID       uint
	BenefitID      uint
	EffectiveDate  string  `gorm:"type:varchar(8)"`
	FundCode       string  `gorm:"type:varchar(5)"` //P0050
	FundType       string  `gorm:"type:varchar(2)"` //P0050
	FundCurr       string  `gorm:"type:varchar(3)"` //P0050
	FundPercentage float64 `gorm:"type:decimal(15,5)"`
	HistoryCode    string  `gorm:"type:varchar(10)"` // Transaciton Code  H0136

}

type IlpSwitchHeader struct {
	gorm.Model
	types.CModel
	PolicyID        uint
	BenefitID       uint
	EffectiveDate   string `gorm:"type:varchar(8)"`
	FundSwitchBasis string `gorm:"type:varchar(1)"` //P0050 Unit or Amount or Percentage
	IlpSwitchFunds  []IlpSwitchFund
}

type IlpSwitchFund struct {
	gorm.Model
	types.CModel
	PolicyID          uint
	BenefitID         uint
	IlpSwitchHeaderID uint
	EffectiveDate     string `gorm:"type:varchar(8)"`
	SwitchDirection   string `gorm:"type:varchar(1)"` //P0050 Source or Target
	SequenceNo        uint
	FundCode          string  `gorm:"type:varchar(5)"`   //P0050
	FundType          string  `gorm:"type:varchar(2)"`   //P0050
	FundCurr          string  `gorm:"type:varchar(3)"`   //P0050
	FundPercentage    float64 `gorm:"type:decimal(8,5)"` // Arrived Value
	FundUnits         float64 `gorm:"type:decimal(15,5);"`
	FundAmount        float64
}

type IlpStatementPrint struct { // temporary table
	gorm.Model
	types.CModel
	PolicyID        uint
	BenefitID       uint
	Description     string  `gorm:"type:varchar(50)"` // opening or closing or ilp prem or ilp fee or ilp contribution
	TransactionDate string  `gorm:"type:varchar(8)"`
	EffectiveDate   string  `gorm:"type:varchar(8)"`
	SequenceNo      uint    // 1,2,3
	FundCode        string  `gorm:"type:varchar(5)"` //P0050
	FundUnits       float64 `gorm:"type:decimal(15,5);"`
	FundAmount      float64
	FundPrince      float64
}

// Read ILPANNSUMMARY TABLE WITH EFFECTIVE AS FROM DATE
// IF RECORD FOUND, MOVE VALUES TO  "opENING"
// READ ILPTRANSACTIONS EFFECITVE DATE >= OPENING DATE AND EFFECTIVE <= CLOSING DATE
// MOVE ILPTRANSACTIONS INTO THIS NEW FILE
// READ ILPANNSUMMARY TABLE WITH EFFECIGE DATE AS CLOSING DATE
// MOVE VALUES TO NEW FILE WITH DESCRIPTION AS 'CLOSING"
