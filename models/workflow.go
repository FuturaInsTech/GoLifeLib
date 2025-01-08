package models

import (
	"time"

	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
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

// type WorkflowPolicy struct {
// 	gorm.Model
// 	types.CModel
// 	PolicyID             uint
// 	WorkflowType         string `gorm:"type:varchar(20)"`
// 	WorkflowTypeRef      string `gorm:"type:varchar(20)"`
// 	WorkflowSeqno        uint
// 	WorkflowDescription  string `gorm:"type:varchar(50)"`
// 	WorkflowSubType      string `gorm:"type:varchar(20)"`
// 	WorkflowSubTypeDesc  string `gorm:"type:varchar(50)"`
// 	WorkflowOptIndicator string `gorm:"type:varchar(1)"` // M Mandatory O Option
// }

type WfTaskAssignment struct {
	gorm.Model
	types.CModel
	TaskID       uint
	AssignedTo   uint
	AssignedAt   time.Time
	StartedAt    time.Time
	CompletedAt  time.Time
	SlaViolation bool   `json:"is_active"`
	Comments     string `gorm:"type:longtext"`
	Attachments  string `gorm:"type:longtext"`
}

type WfActionAssignment struct {
	gorm.Model
	types.CModel
	TaskID       uint
	ActionID     uint
	AssignedTo   uint
	AssignedAt   time.Time
	StartedAt    time.Time
	CompletedAt  time.Time
	SlaViolation bool   `json:"is_active"`
	Comments     string `gorm:"type:longtext"`
	Attachments  string `gorm:"type:longtext"`
}

type WfTaskExecutionLog struct {
	gorm.Model
	types.CModel
	ReqRefData   string `gorm:"type:varchar(20)"`
	ReqRefType   string `gorm:"type:varchar(20)"`
	TaskID       uint
	ActionID     uint
	AssignedTo   uint
	AssignedAt   time.Time
	StartedAt    time.Time
	CompletedAt  time.Time
	SlaViolation bool `json:"is_active"`
	LoggedAt     time.Time
	// action      string `gorm:"type:varchar(10)"` //sub-tasks of tasks ??
	// PerformedBy uint
	Details string `gorm:"type:longtext"`
}

type WfAction struct {
	gorm.Model
	types.CModel
	TaskID               uint
	ActionName           string `gorm:"type:varchar(20)"`
	ActionDescription    string `gorm:"type:longtext"`
	SeqNo                uint
	SlaDuration          uint
	ActionStatus         string `gorm:"type:varchar(2)"`
	Priority             string `gorm:"type:varchar(2)"`
	DueDate              string `gorm:"type:varchar(8)"`
	WfActionAssignmentID uint
	//WfTaskExecutionLogs []WfTaskExecutionLog
}

type WfTask struct {
	gorm.Model
	types.CModel
	RequestID          uint
	TaskName           string `gorm:"type:varchar(20)"`
	TaskDescription    string `gorm:"type:longtext"`
	SeqNo              uint
	SlaDuration        uint
	TaskStatus         string `gorm:"type:varchar(2)"`
	Priority           string `gorm:"type:varchar(2)"`
	DueDate            string `gorm:"type:varchar(8)"`
	WfTaskAssignmentID uint
	//WfActions           []WfAction
	//WfTaskExecutionLogs []WfTaskExecutionLog
}

type WfRequest struct {
	gorm.Model
	types.CModel
	ReqName        string `gorm:"type:varchar(20)"`
	ReqDescription string `gorm:"type:longtext"`
	CreatedBy      uint
	ReqRefData     string `gorm:"type:varchar(20)"`
	ReqRefType     string `gorm:"type:varchar(20)"`
	ReqStatus      string `gorm:"type:varchar(20)"`
	//WfTasks        []WfTask
}

type UserDepartment struct {
	gorm.Model
	types.CModel
	DepartmentID    string `gorm:"type:varchar(20)"`
	TeamsID         string `gorm:"type:varchar(20)"`
	UserID          string `gorm:"type:varchar(20)"`
	UserDesignation string `gorm:"type:varchar(20)"`
	UserLevel       string `gorm:"type:varchar(20)"`
}

type WfUserReminder struct {
	gorm.Model
	types.CModel
	UserId         uint
	ReminderRef    string `gorm:"type:varchar(20)"`
	ReminderOn     string `gorm:"type:varchar(20)"`
	ReminderNote   string `gorm:"type:varchar(20)"`
	ReminderDatime time.Time
	ReminderType   string `gorm:"type:varchar(20)"`
	PhoneNo        string `gorm:"type:varchar(20)"`
	Email          string `gorm:"type:varchar(20)"`
	ReminderPerson string `gorm:"type:varchar(50)"`
}
