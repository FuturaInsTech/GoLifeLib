package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	types.CModel
	Method      string `gorm:"type:varchar(50)"`
	Description string `gorm:"type:varchar(50)"`
	Permissions []Permission
}
