package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	types.CModel
	TranCode    string `gorm:"type:varchar(50)"`
	Method      string `gorm:"type:varchar(50)"`
	Description string `gorm:"type:varchar(50)"`
	Permissions []Permission
}
