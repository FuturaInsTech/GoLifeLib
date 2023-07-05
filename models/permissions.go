package models

import (
	"database/sql"

	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	types.CModel
	ModelName string `gorm:"type:varchar(100)"`
	Method    string `gorm:"type:varchar(100)"`
	// sql.NullInt gives nullable value
	//UserID      sql.NullInt64
	//UserGroupID sql.NullInt64
	UserID        sql.NullInt64
	UserGroupID   sql.NullInt64
	TransactionID uint
}
