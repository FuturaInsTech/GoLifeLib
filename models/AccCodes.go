package models

import (
	"github.com/shijith.chand/go-jwt/types"
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
