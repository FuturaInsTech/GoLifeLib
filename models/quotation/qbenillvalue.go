package quotation

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type QBenIllValue struct {
	gorm.Model
	types.CModel
	QDetailID         uint
	QCoverage         string `gorm:"type:varchar(05)"`
	QSumAssured       uint64
	QMaturityDate     string `gorm:"type:varchar(08)"`
	QMaturityAmt      uint64
	QPolicyYear       uint
	QPolAnnivDate     string `gorm:"type:varchar(08)"`
	QLifeAssuredAge   uint
	QTotalPremPaid    float64
	QDeathBenefitAmt  uint64
	QRevBonusAmt      uint64
	QTerBonusAmt      uint64
	QGuarAdditions    uint64
	QLoyaltyAdditions uint64
	QGuarSurrValue    uint64
	QSplSurrValue     uint64
	QBonusSurValue    uint64
	QAccuDividend     uint64
	QAccuDivInterest  uint64
	QAntiSurBenAmt    uint64 // for holding the Survival Benefit Amount at end of policy year
	QallocatedAmt     float64
	QUnallocedAmt     float64
	QPesValamt        float64
	QNorValamt        float64
	QOptValamt        float64
	QPesMCamt         float64
	QNorMCamt         float64
	QOptMCamt         float64
	QPesFEamt         float64
	QNorFEamt         float64
	QOptFEamt         float64
}
