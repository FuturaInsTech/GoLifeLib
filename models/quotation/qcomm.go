package quotation

import (
	"github.com/shijith.chand/go-jwt/types"
	"gorm.io/gorm"
)

type Qcommunication struct {
	gorm.Model
	types.CModel
	TemplateName         string `gorm:"type:varchar(10)"` //P0033
	Language             string `gorm:"type:varchar(02)"` //P0002
	QHeaderID            uint
	Seqno                uint16 `gorm:"asc"`
	ClientID             uint
	AgencyID             uint
	EffectiveDate        string          `gorm:"type:varchar(08)"`
	ExtractedDate        string          `gorm:"type:varchar(08)"`
	ExtractedStaus       string          `gorm:"type:varchar(02)"`
	ExtractedData        types.ExtraData `gorm:"type:text(99999)"`
	SMSAllowed           string          `gorm:"type:varchar(01)"`
	EmailAllowed         string          `gorm:"type:varchar(01)"`
	WhatsAppAllowed      string          `gorm:"type:varchar(01)"`
	AgentSMSAllowed      string          `gorm:"type:varchar(01)"`
	AgentEmailAllowed    string          `gorm:"type:varchar(01)"`
	AgentWhatsAppAllowed string          `gorm:"type:varchar(01)"`
	CompanyEmail         string          `gorm:"type:varchar(80)"`
	CompanyPhone         string          `gorm:"type:varchar(20)"`
	DepartmentName       string          `gorm:"type:varchar(50)"`
	DepartmentHead       string          `gorm:"type:varchar(50)"`
	TemplatePath         string          `gorm:"type:varchar(50)"`
	PDFPath              string          `gorm:"type:varchar(50)"`
}
