package models

import (
	"time"

	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type TransactionLock struct {
	CompanyID     uint             `gorm:"primaryKey"`
	LockedType    types.LockedType `gorm:"type:tinyint unsigned;primaryKey"`
	LockedTypeKey string           `gorm:"type:varchar(30);primaryKey;"`
	VersionId     string           `gorm:"type:varchar(40)"`
	IsValid       types.BitBool    `gorm:"type:bit(1)"`
	IsLocked      types.BitBool    `gorm:"type:bit(1)"`
	UpdatedID     uint64
	Tranno        uint
	Session       string `gorm:"type:varchar(15)"`
	SessionID     uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
