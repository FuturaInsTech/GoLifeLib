package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type TDFParam struct {
	gorm.Model
	types.CModel
	FromPolicy uint
	ToPolicy   uint
}
