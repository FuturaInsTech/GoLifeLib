package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type FieldValidator struct {
	gorm.Model
	types.CModel
	FunctionName     string `gorm:"type:varchar(50)"`
	FieldName        string `gorm:"type:varchar(50)"`
	ParamName        string `gorm:"type:varchar(50)"`
	LanguageId       uint
	CompanyID        uint
	ErrorDescription string `gorm:"type:varchar(200)"`
	BlankAllowed     string `gorm:"type:varchar(01)"`
	ZeroAllowed      string `gorm:"type:varchar(01)"`
}
