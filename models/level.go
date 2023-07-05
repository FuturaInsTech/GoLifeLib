package models

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type Level struct {
	gorm.Model
	types.CModel
	ShortCode string `gorm:"type:varchar(05)"`
	LongName  string `gorm:"type:varchar(40)"`
	LevelCode string `gorm:"type:varchar(03)"`
	Levels    []Level
	LevelId   uint64
}
