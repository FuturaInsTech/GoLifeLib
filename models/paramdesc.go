package models

import (
	"time"
)

type ParamDesc struct {
	CompanyId   uint16 `gorm:"primaryKey;"`
	Name        string `gorm:"type:varchar(50);primaryKey;"`
	Item        string `gorm:"type:varchar(50);primaryKey;"`
	RecType     string `gorm:"type:varchar(2);primaryKey;"`
	LanguageId  uint8  `gorm:"primaryKey;"`
	Shortdesc   string `gorm:"type:varchar(20);"`
	Longdesc    string `gorm:"type:varchar(50);"`
	LastModUser uint64 `gorm:"type:bigint;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
