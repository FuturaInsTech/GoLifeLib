package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type TDFParam struct {
	gorm.Model
	types.CModel
	FromPolicy uint
	ToPolicy   uint
}
