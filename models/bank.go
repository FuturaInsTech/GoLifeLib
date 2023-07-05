package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type Bank struct {
	gorm.Model
	types.CModel
	BankCode          string `gorm:"type:varchar(50)"` //IFSC
	BankAccountNo     string `gorm:"type:varchar(50)"` // Bank Account No of Payer/Agent/Claiment
	StartDate         string `gorm:"type:varchar(08)"`
	EndDate           string `gorm:"type:varchar(08)"`
	BankType          string `gorm:"type:varchar(05)"` //P0020
	BankAccountStatus string `gorm:"type:varchar(05)"` //P0021
	CompanyID         uint
	ClientID          uint
	Payers            []Payer
	Agencies          []Agency
	MedProviders      []MedProvider
}
