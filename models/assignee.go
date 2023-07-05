package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Assignee struct {
	gorm.Model
	types.CModel
	PolicyID     uint
	ClientID     uint
	AssigneeType string `gorm:"type:varchar(20)"` // P0052
	FromDate     string `gorm:"type:varchar(08)"`
	ToDate       string `gorm:"type:varchar(08)"`
}
