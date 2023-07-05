package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Nominee struct {
	gorm.Model
	types.CModel
	PolicyID            uint
	ClientID            uint
	NomineeRelationship string `gorm:"type:varchar(05)"` //P0045
	NomineeLongName     string `gorm:"type:varchar(100)"`
	NomineePercentage   float64
}
