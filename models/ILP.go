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
	FundCurrency   string `gorm:"type:varchar(3)"`
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
	FundCode  string `gorm:"type:varchar(5)"` //P0050
	FundType  string `gorm:"type:varchar(2)"` //ACcummulated/INinitial/BOus P0050
	FundUnits float64
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
	FundCurrency        string `gorm:"type:varchar(3)"`
	FundUnits           float64
	FundPrice           float64
	CurrentOrFuture     string `gorm:"type:varchar(1)"` //P0050
	OriginalAmount      float64
	OriginalCurrency    string `gorm:"type:varchar(3)"`
	HistoryCode         string `gorm:"type:varchar(5)"`
	InvNonInvFlag       string `gorm:"type:varchar(2)"`
	InvNonInvPercentage float64
	AccountCodeID       uint
	AccountCode         string `gorm:"type:varchar(30)"`
	CurrencyRate        float64
	MortalityIndicator  string `gorm:"type:varchar(1)"`
	SurrenderPercentage float64
}
