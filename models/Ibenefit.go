package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type IBenefit struct {
	gorm.Model
	types.CModel
	PolicyID                     uint
	BenefitID                    uint
	BCoverage                    string `gorm:"type:varchar(05)"` //Q0006
	Seqno                        uint
	PayFrequency                 string `gorm:"type:varchar(1)"` //10
	Percentage                   float64
	BSumAssured                  uint64
	EffectiveDate                string  `gorm:"type:varchar(8)"`
	IncidentDate                 string  `gorm:"type:varchar(8)"`
	ReceivedDate                 string  `gorm:"type:varchar(8)"`
	PaidToDate                   string  `gorm:"type:varchar(8)"`
	BStatusCode                  string  `gorm:"type:varchar(2)"`
	ApprovalFlag                 string  `gorm:"type:varchar(2)"`
	CertificateExistranceFlag    string  `gorm:"type:varchar(1)"` // P0058
	CertificateExistranceDate    string  `gorm:"type:varchar(8)"` // 20250101
	CertificateExistranceRevDate string  `gorm:"type:varchar(8)"` // 20250131
	NextPayDate                  string  `gorm:"type:varchar(8)"`
	ClaimAmount                  float64 //  20000
	BStartDate                   string  `gorm:"type:varchar(8)"`
	TotalAmount                  float64
	PaidDate                     string `gorm:"type:varchar(8)"`
}
