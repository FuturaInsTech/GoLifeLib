package models

import (
	"github.com/FuturaInsTech/GoLifeLib/models/quotation"
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Client struct {
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
	ClientPaID        uint
	Addresses         []Address
	Nominees          []Nominee
	LeadDetails       []LeadDetail
	QHeaders          []quotation.QHeader
	Agencies          []Agency
	Banks             []Bank
	Policies          []Policy
	Owners            []Payer
	Benefits          []Benefit
	Receipts          []Receipt
	Communications    []Communication
	MedProviders      []MedProvider
	MedReqs           []MedReq
	QDetails          []quotation.QDetail
	LeadAllocations   []LeadAllocation
	LeadFollowups     []LeadFollowup
	DeathHs           []DeathH
	DeathDs           []DeathD
	Payers            []Payer
	SaChanges         []SaChange
	Mrtas             []Mrta
	SurrHs            []SurrH
	SurrDs            []SurrD
	Payments          []Payment
	PayingAuthorities []PayingAuthority
}

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
