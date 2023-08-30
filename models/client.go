package models

import (
	"github.com/FuturaInsTech/GoLifeLib/models/quotation"
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	types.CModel

	ClientShortName string `gorm:"type:varchar(50)"`
	ClientLongName  string `gorm:"type:varchar(50)"`
	ClientSurName   string `gorm:"type:varchar(50)"`
	Gender          string `gorm:"type:varchar(05)"`
	Salutation      string `gorm:"type:varchar(05)"`
	Language        string `gorm:"type:varchar(05)"`
	ClientDob       string `gorm:"type:varchar(8)"`
	ClientDod       string `gorm:"type:varchar(8)"`
	ClientEmail     string `gorm:"type:varchar(100)"`
	ClientMobile    string `gorm:"type:varchar(20)"`
	ClientStatus    string `gorm:"type:varchar(05)"`
	ClientType      string `gorm:"type:varchar(01)"` // C CORPORATE I FOR INDIVIDUAL
	Addresses       []Address
	Nominees        []Nominee
	LeadDetails     []LeadDetail
	QHeaders        []quotation.QHeader
	Agencies        []Agency
	Banks           []Bank
	Policies        []Policy
	Owners          []Payer
	Benefits        []Benefit
	Receipts        []Receipt
	Communications  []Communication
	MedProviders    []MedProvider
	MedReqs         []MedReq
	QDetails        []quotation.QDetail
	LeadAllocations []LeadAllocation
	LeadFollowups   []LeadFollowup
	DeathHs         []DeathH
	DeathDs         []DeathD
	Payers          []Payer
	SaChanges       []SaChange
	Mrtas           []Mrta
	SurrHs          []SurrH
	SurrDs          []SurrD
	Payments        []Payment
}
