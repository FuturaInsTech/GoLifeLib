package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type Extra struct {
	gorm.Model
	types.CModel
	PolicyID    uint
	EReason     string `gorm:"type:varchar(02)"` //P0026
	EMethod     string `gorm:"type:varchar(05)"` //P0025
	EPrem       float64
	EPercentage float64
	EAmt        float64
	ETerm       int
	EAge        int
	BenefitID   uint
	BCoverage   string `gorm:"type:varchar(05)"`
	FromDate    string `gorm:"type:varchar(08)"`
	ToDate      string `gorm:"type:varchar(08)"`
}
