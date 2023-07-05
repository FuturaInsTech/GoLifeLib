package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type MedProvider struct {
	gorm.Model
	types.CModel
	MedProviderName   string `gorm:"type:varchar(100)"`
	ClientID  uint
	AddressID uint 
	BankID  	uint 
	MedCurr    string `gorm:"type:varchar(03)"`  // P0023 
	StartDate  string `gorm:"type:varchar(08)"`
	EndDate    string `gorm:"type:varchar(08)"`
	MedStatus  string `gorm:"type:varchar(03)"` //P0004 
	MedReqs     [] MedReq	

}