package models

import (
	"database/sql"

	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type PHistory struct {
	gorm.Model
	types.CModel
	HistoryCode   string `gorm:"type:varchar(05)"`
	PolicyID      uint
	EffectiveDate string `gorm:"type:varchar(08)"`
	CurrentDate   string `gorm:"type:varchar(08)"`
	ReversedAt    sql.NullTime
	Is_reversed   types.BitBool `gorm:"type:bit(1)"`
	PrevData      types.ExtraData
	RevRemark     string `gorm:"type:varchar(500)"`
	RevUserID     sql.NullInt64
}
