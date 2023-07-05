package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Gender struct {
	gorm.Model
	types.CModel
	GenderShortName string `gorm:"type:varchar(1)"`
	GenderLongName  string `gorm:"type:varchar(10)"`
}
