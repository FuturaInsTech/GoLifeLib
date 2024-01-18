package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type ClientPa struct {
	gorm.Model
	types.CModel

	ClientID          uint
	PayingAuthorityID uint   // Employed with Pa Reference
	PayRollNumber     string `gorm:"type:varchar(20)"` // Employee Pay Roll Number
	Designation       string `gorm:"type:varchar(30)"` // Employee Designation
	Department        string `gorm:"type:varchar(20)"` // Employed in Department
	Location          string `gorm:"type:varchar(20)"` // Employed at Location
	StartDate         string `gorm:"type:varchar(08)"` // Employment Start Date
	EndDate           string `gorm:"type:varchar(08)"` // Employment End Date
	PrevPaID          uint   // Use when change in Employment Only
}
