package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type UserGroup struct {
	gorm.Model
	types.CModel
	GroupName   string `gorm:"type:varchar(100)"`
	ValidFrom   string `gorm:"type:varchar(08)"`
	ValidTo     string `gorm:"type:varchar(08)"`
	Users       []User
	Permissions []Permission
}
