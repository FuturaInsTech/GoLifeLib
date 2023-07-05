package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Error struct {
	gorm.Model
	types.CModel
	ShortCode  string `gorm:"type:varchar(05)"`
	LongCode   string `gorm:"type:varchar(80)"`
	LanguageID uint
}
