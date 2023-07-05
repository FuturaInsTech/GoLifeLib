package models

import "gorm.io/gorm"

type Relationship struct {
	gorm.Model
	RelationshipShortName string `gorm:"type:varchar(08)"`
	RelationshipLongName  string `gorm:"type:varchar(100)"`
}
