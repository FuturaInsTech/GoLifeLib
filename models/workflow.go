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
