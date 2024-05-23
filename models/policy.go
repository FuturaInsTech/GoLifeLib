package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type Policy struct {
	gorm.Model
	types.CModel
	PRCD            string `gorm:"type:varchar(08)"`
	ProposalDate    string `gorm:"type:varchar(08)"`
	PProduct        string `gorm:"type:varchar(05)"` //Q0005
	PFreq           string `gorm:"type:varchar(01)"` //Q0009
	PContractCurr   string `gorm:"type:varchar(03)"` //P0023
	PBillCurr       string `gorm:"type:varchar(03)"` //P0023
	POffice         string `gorm:"type:varchar(05)"` //P0018
	PolStatus       string `gorm:"type:varchar(05)"` //P0024 READ ONLY
	PReceivedDate   string `gorm:"type:varchar(08)"`
	PUWDate         string `gorm:"type:varchar(08)"`
	ClientID        uint   // Owner Client ID
	BTDate          string `gorm:"type:varchar(08)"` // READ ONLY
	PaidToDate      string `gorm:"type:varchar(08)"` // READONLY
	NxtBTDate       string `gorm:"type:varchar(08)"` // READONLY
	AnnivDate       string `gorm:"type:varchar(08)"` // READONLY
	AgencyID        uint   // NEED TO SELECT
	InstalmentPrem  float64
	BillingType     string `gorm:"type:varchar(05)"` // P0055
	BankID          uint
	PayingAuthority uint   // Paying Authority Client
	NfoMethod       string `gorm:"type:varchar(05)"` //q0005
	Benefits        []Benefit
	PHistories      []PHistory
	Extras          []Extra
	//Receipts        []Receipt
	TDFpolicies []TDFPolicy
	SurvBs      []SurvB
	//Communications   []Communication
	AddressID        uint //Api. GetAllAddressByClientID
	Uwreasons        []Uwreason
	MedReqs          []MedReq
	Payers           []Payer
	Nominees         []Nominee
	DeathHs          []DeathH
	DeathDs          []DeathD
	SaChanges        []SaChange
	Mrtas            []Mrta
	SurrHs           []SurrH
	SurrDs           []SurrD
	Tdfhs            []Tdfh
	TDFPolicies      []TDFPolicy
	MaturityHs       []MaturityH
	MaturityDs       []MaturityD
	PolBills         []PolBill
	CriticalIllnesss []CriticalIllness
	IBenefits        []IBenefit
	Payments         []Payment
	IlpFunds         []IlpFund
	IlpTransactions  []IlpTransaction
	IlpSummaries     []IlpSummary
	IlpSwitchHeaders []IlpSwitchHeader
	IlpSwitchFunds   []IlpSwitchFund
}
