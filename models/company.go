package models

import (
	"github.com/FuturaInsTech/GoLifeLib/models/quotation"
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	CompanyName              string `gorm:"type:varchar(80)"`
	CompanyAddress1          string `gorm:"type:varchar(80)"`
	CompanyAddress2          string `gorm:"type:varchar(80)"`
	CompanyAddress3          string `gorm:"type:varchar(80)"`
	CompanyAddress4          string `gorm:"type:varchar(80)"`
	CompanyAddress5          string `gorm:"type:varchar(80)"`
	CompanyPostalCode        string `gorm:"type:varchar(80)"`
	CompanyCountry           string `gorm:"type:varchar(80)"`
	CompanyUid               string `gorm:"type:varchar(40)"`
	CompanyGst               string `gorm:"type:varchar(40)"`
	CompanyPan               string `gorm:"type:varchar(40)"`
	CompanyTan               string `gorm:"type:varchar(40)"`
	CompanyLogo              string `gorm:"type:longtext"`
	CompanyIncorporationDate string `gorm:"type:varchar(08)"`
	CompanyTerminationDate   string `gorm:"type:varchar(08)"`
	CompanyStatusID          uint
	CurrencyID               uint   // P0030  USD2INR
	NationalIdentityMand     string `gorm:"type:varchar(01)"`
	NationalIdentityEncrypt  string `gorm:"type:varchar(01)"`

	ContHeaders     []ContHeader
	Users           []User
	LeadChannels    []LeadChannel
	LeadDetails     []LeadDetail
	Errors          []Error
	Campaigns       []Campaign
	CampaignComps   []CampaignComp
	LeadAllocations []LeadAllocation
	Clients         []Client
	Addresses       []Address
	QHeaders        []quotation.QHeader
	QDetails        []quotation.QDetail
	UserGroups      []UserGroup
	// Agencies        []Agency
	FieldValidators []FieldValidator
	Levels          []Level
	QBenIllValues   []quotation.QBenIllValue
	Banks           []Bank
	Permissions     []Permission
	Extras          []Extra
	Policies        []Policy
	Benefits        []Benefit
	PHistories      []PHistory
	GlTypes         []GlType
	AccCodes        []AccountCode
	GlBals          []GlBal
	Receipts        []Receipt
	Tdfhs           []Tdfh
	TDFRules        []TDFRule
	TDFPolicies     []TDFPolicy
	Communications  []Communication
	LeadFollowups   []LeadFollowup
	//CreatedAt time.Time
	//UpdatedAt time.Time
	Transactions      []Transaction
	Uwreasons         []Uwreason
	MedProviders      []MedProvider
	MedReqs           []MedReq
	Payers            []Payer
	SaChanges         []SaChange
	SurrHs            []SurrH
	SurrDs            []SurrD
	BusinessDates     []BusinessDate
	MaturityHs        []MaturityH
	MaturityDs        []MaturityD
	PolBills          []PolBill
	PayingAuthorities []PayingAuthority
	PaBillSummaries   []PaBillSummary
}
