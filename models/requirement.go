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
	Remarks      string `gorm:"type:varchar(160)"`
	InvoieNo     uint   // Update Invoice No from ReqBill
	InvoiceDate  string `gorm:"type:varchar(08)"` // Update Invoice Date from ReqBill
}

type ReqProof struct {
	gorm.Model
	types.CModel
	ReqcallID uint
	ProofImg  string `gorm:"type:longtext"`
}

type ReqBill struct {
	gorm.Model
	types.CModel
	MedId            uint
	InvoieNo         uint
	InvoiceDate      string `gorm:"type:varchar(08)"`
	InvoiceClient    uint
	InvoicePolicy    uint
	ExaminationDate  string `gorm:"type:varchar(08)"`
	InvoiceAmount    float64
	InvoiceRemarks   string `gorm:"type:varchar(160)"`
	CreatedDate      string `gorm:"type:varchar(08)"`
	ReconcileStatus  string `gorm:"type:varchar(02)"`
	ReconcileDate    string `gorm:"type:varchar(08)"`
	ReconcileAmount  float64
	ReconcileRemarks string `gorm:"type:varchar(160)"`
	PayType          string `gorm:"type:varchar(02)"`
	PayFlag          string `gorm:"type:varchar(02)"`
	PayDate          string `gorm:"type:varchar(08)"`
	PayReference     uint
}

type UserLimit struct {
	gorm.Model
	types.CModel
	UserId      uint
	HistoryCode string `gorm:"type:varchar(05)"`
	From        float64
	To          float64
}
