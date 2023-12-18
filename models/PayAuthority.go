package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type PayingAuthority struct {
	gorm.Model
	types.CModel
	ClientID uint
	Policies []Policy
}
