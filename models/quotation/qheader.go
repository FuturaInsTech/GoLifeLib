package quotation

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type QHeader struct {
	gorm.Model
	types.CModel
	QuoteDate     string `gorm:"type:varchar(08)"`
	POffice       string `gorm:"type:varchar(05)"` //P0018
	QProduct      string `gorm:"type:varchar(05)"`
	QContractCurr string `gorm:"type:varchar(03)"`
	ClientID      uint
	QFirstName    string `gorm:"type:varchar(30)"`
	QLastName     string `gorm:"type:varchar(30)"`
	QMidName      string `gorm:"type:varchar(30)"`
	QDob          string `gorm:"type:varchar(08)"`
	QGender       string `gorm:"type:varchar(01)"`
	QNri          string `gorm:"type:varchar(01)"`
	QEmail        string `gorm:"type:varchar(100)"`
	QMobile       string `gorm:"type:varchar(20)"`
	QStatus       string `gorm:"type:varchar(10)"`
	// Q0007
	QOccGroup string `gorm:"type:varchar(04)"`
	// Q0008
	QOccSect        string `gorm:"type:varchar(04)"`
	QOccupation     string `gorm:"type:varchar(30)"`
	QAnnualIncome   uint64
	QDeclaration    string `gorm:"type:varchar(30)"`
	AddressID       uint
	AgencyID        uint
	QDetails        []QDetail
	Finalized       string `gorm:"type:varchar(01)"`
	Qcommunications []Qcommunication
}
