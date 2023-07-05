package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Uwreason struct {
	gorm.Model
	types.CModel
	ReasonDescription string `gorm:"type:varchar(1000)"` // Mandatory
	PolicyID          uint
	RequestedDate     string `gorm:"type:varchar(08)"` // Mandatory
}
