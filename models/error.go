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

type TxnError struct {
	ErrorCode string `gorm:"type:varchar(05)"`
	ParamName string `gorm:"type:varchar(20)"`
	ParamItem string `gorm:"type:varchar(20)"`
	DB        string `gorm:"type:varchar(20)"`
	DbError   string `gorm:"type:varchar(50)"`
}
