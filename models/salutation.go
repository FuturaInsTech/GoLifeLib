package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type Salutation struct {
	gorm.Model
	types.CModel
	ShortName string `gorm:"type:varchar(8)"`
	LongName  string `gorm:"type:varchar(50)"`
}
