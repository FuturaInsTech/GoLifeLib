package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type Payer struct {
	gorm.Model
	types.CModel
	PolicyID uint
	ClientID uint
	BankID   uint
	FromDate string `gorm:"type:varchar(08)"`
	ToDate   string `gorm:"type:varchar(08)"`
}
