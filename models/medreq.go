package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type MedReq struct{
	gorm.Model
	types.CModel
	CreatedDate  	string `gorm:"type:varchar(08)"`
	EffectiveDate 	string `gorm:"type:varchar(08)"`
	ReminderDate  	string `gorm:"type:varchar(08)"`
	PolicyID 		uint
	BenefitID 		uint
	MedCode   		string `gorm:"type:varchar(08)"` //P0040
	Seqno           int 
	Status          string `gorm:"type:varchar(01)"` //P0042
	MedProviderID   uint 
	ClientID        uint 
}