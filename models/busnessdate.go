package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type BusinessDate struct {
	gorm.Model
	types.CModel
	UserID     uint
	Department string `gorm:"type:varchar(05)"`
	Date       string `gorm:"type:varchar(08)"`
}
