package models

import (
	"time"

	"github.com/FuturaInsTech/GoLifeLib/types"
)

type Param struct {
	CompanyId   uint16          `gorm:"primaryKey;"`
	Name        string          `gorm:"type:varchar(50);primaryKey;"`
	Item        string          `gorm:"type:varchar(50);primaryKey;"`
	RecType     string          `gorm:"type:varchar(2);primaryKey;"`
	Seqno       uint16          `gorm:"primaryKey"`
	StartDate   string          `gorm:"type:varchar(8)"`
	EndDate     string          `gorm:"type:varchar(8)"`
	Is_valid    types.BitBool   `gorm:"type:bit(1)"`
	Data        types.ExtraData `gorm:"type:varchar(5000)"`
	LastModUser uint64          `gorm:"type:bigint;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
