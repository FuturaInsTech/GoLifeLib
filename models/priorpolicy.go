package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type PriorPolicy struct {
	gorm.Model
	types.CModel
	PolicyID         uint
	SeqNo            uint
	PriorPolNo       string `gorm:"type:varchar(50)"`
	PriorInsurerName string `gorm:"type:varchar(100)"`
	PStartDate       string `gorm:"type:varchar(08)"`
	PEndDate         string `gorm:"type:varchar(08)"`
	PSumAssured      uint64
}
