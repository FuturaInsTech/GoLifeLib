package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type ReqCall struct { //RequirementCall
	gorm.Model
	types.CModel
	PolicyID uint
	ClientID uint
	// AssigneeType string `gorm:"type:varchar(20)"` // P0052
	CreateDate   string `gorm:"type:varchar(08)"`
	ReqType      string `gorm:"type:varchar(20)"` // P0050
	ReqCode      string `gorm:"type:varchar(20)"` // P0050
	DocDate      string `gorm:"type:varchar(08)"`
	CompleteDate string `gorm:"type:varchar(08)"`
	RemindDate   string `gorm:"type:varchar(08)"`
	ReqStatus    string `gorm:"type:varchar(20)"` // P0050
	MedId        uint
	PayClientID  uint
}

type ReqProof struct {
	gorm.Model
	types.CModel
	ReqcallID uint
	ProofImg  string `gorm:"type:longtext"`
}
