package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type PayingAuthority struct {
	gorm.Model
	types.CModel
	ClientID  uint
	PaName    string
	PaType    string //P0050
	StartDate string
	EndDate   string
	PaStatus  string // P0050
}
