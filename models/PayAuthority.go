package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type PayingAuthority struct {
	gorm.Model
	types.CModel
	ClientID  uint
	PaName    string `gorm:"type:varchar(50)"`
	PaType    string `gorm:"type:varchar(01)"` //P0050
	StartDate string `gorm:"type:varchar(08)"`
	EndDate   string `gorm:"type:varchar(08)"`
	PaStatus  string `gorm:"type:varchar(02)"` // P0050
}

type PaBillSummary struct {
	gorm.Model
	types.CModel
	PaID               uint
	PaBillDueMonth     string `gorm:"type:varchar(06)"` // YYYYMM
	PaBillSeqNo        uint
	PaBillStatus       string `gorm:"type:varchar(01)"` // P0050
	ExtractedDate      string `gorm:"type:varchar(8)"`
	ExtractedCount     uint
	ExtractedAmount    float64
	DeductedCount      uint
	DeductedAmount     float64
	NotDeductedCount   uint
	NotDeductedAmount  float64
	UnReconciledCount  uint
	UnReconciledAmount float64
	ReconciledDate     string `gorm:"type:varchar(8)"`
	ReconciledBy       string `gorm:"type:varchar(30)"`
	ApprovedDate       string `gorm:"type:varchar(8)"`
	ApprovedBy         string `gorm:"type:varchar(30)"`
}
