package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Bank struct {
	gorm.Model
	types.CModel
	BankCode          string `gorm:"type:varchar(50)"` //IFSC
	BankAccountNo     string `gorm:"type:varchar(50)"` // Bank Account No of Payer/Agent/Claiment
	StartDate         string `gorm:"type:varchar(08)"`
	EndDate           string `gorm:"type:varchar(08)"`
	BankType          string `gorm:"type:varchar(05)"` //P0020 Savings/Current Bank A/C
	BankAccountStatus string `gorm:"type:varchar(05)"` //P0021
	BankGroup         string `gorm:"type:varchar(05)"` //P0056
	CompanyID         uint
	ClientID          uint
	Payers            []Payer
	Agencies          []Agency
	MedProviders      []MedProvider
}
