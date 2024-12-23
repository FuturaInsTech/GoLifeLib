package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

type WorkflowComments struct {
	gorm.Model
	types.CModel
	PolicyID        uint
	ClientID        uint
	UserID          uint
	Sequence        uint
	WorkflowComment string `gorm:"type:varchar(2500)"`
	CommentProof    string `gorm:"type:longtext"`
}

type WorkflowRules struct {
	types.CModel
	WorkflowSeqno        uint
	WorkflowType         string `gorm:"type:varchar(20)"`
	WorkflowDescription  string `gorm:"type:varchar(50)"`
	WorkflowSubType      string `gorm:"type:varchar(20)"`
	WorkflowSubTypeDesc  string `gorm:"type:varchar(50)"`
	WorkflowOptIndicator string `gorm:"type:varchar(1)"` // M Mandatory O Option
}
