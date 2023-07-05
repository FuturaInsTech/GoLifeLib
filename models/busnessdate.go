package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type BusinessDate struct {
	gorm.Model
	types.CModel
	UserID     uint
	Department string `gorm:"type:varchar(05)"`
	Date       string `gorm:"type:varchar(08)"`
}
