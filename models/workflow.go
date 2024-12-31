package models

import (
	"github.com/FuturaInsTech/GoLifeLib/types"
)

// type WorkflowComments struct {
// 	gorm.Model
// 	types.CModel
// 	PolicyID        uint
// 	ClientID        uint
// 	UserID          uint
// 	Sequence        uint
// 	WorkflowComment string `gorm:"type:varchar(2500)"`
// 	CommentProof    string `gorm:"type:longtext"`
// }

type WorkflowRules struct {
	types.CModel
	WorkflowType         string `gorm:"primaryKey;type:varchar(20)"`
	WorkflowSeqno        uint   `gorm:"primaryKey;"`
	WorkflowDescription  string `gorm:"type:varchar(50)"`
	WorkflowSubType      string `gorm:"type:varchar(20)"`
	WorkflowSubTypeDesc  string `gorm:"type:varchar(50)"`
	WorkflowOptIndicator string `gorm:"type:varchar(1)"` // M Mandatory O Option
}
