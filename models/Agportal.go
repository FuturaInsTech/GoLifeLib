package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type ApClient struct {
	gorm.Model
	types.CModel

	ClientShortName   string `gorm:"type:varchar(50)"`
	ClientLongName    string `gorm:"type:varchar(50)"`
	ClientSurName     string `gorm:"type:varchar(50)"`
	Gender            string `gorm:"type:varchar(05)"`
	Salutation        string `gorm:"type:varchar(05)"`
	Language          string `gorm:"type:varchar(05)"`
	ClientDob         string `gorm:"type:varchar(8)"`
	ClientDod         string `gorm:"type:varchar(8)"`
	ClientEmail       string `gorm:"type:varchar(100)"`
	ClientMobCode     string `gorm:"type:varchar(05)"`
	ClientMobile      string `gorm:"type:varchar(20)"`
	ClientStatus      string `gorm:"type:varchar(05)"`
	ClientType        string `gorm:"type:varchar(01)"` // C CORPORATE I FOR INDIVIDUAL
	NationalId        string `gorm:"type:varchar(50);unique"`
	Nationality       string `gorm:"type:varchar(02)"`
	ClientAltEmail    string `gorm:"type:varchar(100)"`
	ClientAltMobCode  string `gorm:"type:varchar(05)"`
	ClientAltMobile   string `gorm:"type:varchar(20)"`
	ClientWorkID      uint
	CustomerPortal    string `gorm:"type:varchar(01)"` // Yes or No
	CusomterDnd       string `gorm:"type:varchar(01)"` // Yes or No
	ClientReferenceId uint
}

type ApAddress struct {
	gorm.Model
	types.CModel
	AddressType        string `gorm:"type:varchar(05)"` //P0022
	AddressLine1       string `gorm:"type:varchar(50)"`
	AddressLine2       string `gorm:"type:varchar(50)"`
	AddressLine3       string `gorm:"type:varchar(50)"`
	AddressLine4       string `gorm:"type:varchar(50)"`
	AddressLine5       string `gorm:"type:varchar(50)"`
	AddressPostCode    string `gorm:"type:varchar(10)"`
	AddressState       string `gorm:"type:varchar(50)"`
	AddressCountry     string `gorm:"type:varchar(50)"`
	AddressStartDate   string `gorm:"type:varchar(8)"`
	AddressEndDate     string `gorm:"type:varchar(8)"`
	ClientID           uint
	AddressReferenceId uint
}
