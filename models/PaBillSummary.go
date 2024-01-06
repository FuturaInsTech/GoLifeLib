package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type PaBillSummary struct {
	gorm.Model
	types.CModel
	PayingAuthorityID  uint
	PaName             string `gorm:"type:varchar(50)"`
	PaType             string `gorm:"type:varchar(01)"`
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
