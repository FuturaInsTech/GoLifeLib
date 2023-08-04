package quotation

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type QDetail struct {
	gorm.Model
	types.CModel
	QHeaderID uint
	// Q0006
	QCoverage      string `gorm:"type:varchar(05)"`
	QDate          string `gorm:"type:varchar(08)"`
	QRiskSeqNo     uint   // for managing inbuilt multiple covers
	QAge           uint
	QSumAssured    uint64
	QRiskCessAge   uint
	QRiskCessTerm  uint
	QRiskCessDate  string `gorm:"type:varchar(08)"`
	QPremCessAge   uint
	QPremCessTerm  uint
	QPremCessDate  string `gorm:"type:varchar(08)"`
	QBeneCessAge   uint
	QBeneCessTerm  uint
	QBeneCessDate  string `gorm:"type:varchar(08)"`
	QAnnualPremium float64
	QHlyPrem       float64
	QQlyPrem       float64
	QMlyPrem       float64
	QEmrRating     uint
	QAgeAdmitted   string `gorm:"type:varchar(01)"`
	ClientID       uint
	Interest       float64
	QBenIllValues  []QBenIllValue
}
