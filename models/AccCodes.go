package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type AccountCode struct {
	gorm.Model
	types.CModel
	AccountCode string `gorm:"type:varchar(30)"`
	GlSign      string `gorm:"type:varchar(01)"` // + or -
	GlTypeID    uint
	GlMoves     []GlMove
}
