package models

import (
	"time"

	"github.com/FuturaInsTech/GoLifeLib/types"
	"gorm.io/gorm"
)

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
	AssignedFrom uint
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
	AssignedFrom uint
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
	ActionName           string `gorm:"type:varchar(100)"`
	ActionDescription    string `gorm:"type:longtext"`
	TranCode             string `gorm:"type:varchar(10)"`
	SeqNo                uint
	SlaDuration          uint
	SlaDurationType      string `gorm:"type:varchar(2)"`
	ActionStatus         string `gorm:"type:varchar(2)"`
	UpdatedStatusAt      time.Time
	Priority             string `gorm:"type:varchar(2)"`
	DueDate              string `gorm:"type:varchar(8)"`
	DepartmentCode       string `gorm:"type:varchar(20)"`
	TeamsCode            string `gorm:"type:varchar(20)"`
	WfActionAssignmentID uint
	//WfTaskExecutionLogs []WfTaskExecutionLog
}

type WfTask struct {
	gorm.Model
	types.CModel
	RequestID          uint
	TaskName           string `gorm:"type:varchar(100)"`
	TaskDescription    string `gorm:"type:longtext"`
	SeqNo              uint
	SlaDuration        uint
	SlaDurationType    string `gorm:"type:varchar(2)"`
	TaskStatus         string `gorm:"type:varchar(2)"`
	UpdatedStatusAt    time.Time
	Priority           string `gorm:"type:varchar(2)"`
	DueDate            string `gorm:"type:varchar(8)"`
	DepartmentCode     string `gorm:"type:varchar(20)"`
	TeamsCode          string `gorm:"type:varchar(20)"`
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
	ReqTokenNo     string `gorm:"type:varchar(20)"`
	ReqRefData     string `gorm:"type:varchar(20)"`
	ReqRefType     string `gorm:"type:varchar(20)"`
	ReqStatus      string `gorm:"type:varchar(20)"`
	//WfTasks        []WfTask
}

type UserDepartment struct {
	gorm.Model
	types.CModel
	DepartmentCode  string `gorm:"type:varchar(20)"`
	TeamsCode       string `gorm:"type:varchar(20)"`
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

// View Models
type TaskView struct {
	ID              uint      `gorm:"column:id"`
	CompanyID       uint      `gorm:"column:company_id"`
	RequestID       uint      `gorm:"column:request_id"`
	ReqName         string    `gorm:"column:req_name"`
	TaskName        string    `gorm:"column:task_name"`
	TaskDescription string    `gorm:"column:task_description"`
	SeqNo           uint      `gorm:"column:seq_no"`
	SlaDuration     uint      `gorm:"column:sla_duration"`
	SlaDurationType string    `gorm:"column:sla_duration_type"`
	TaskStatus      string    `gorm:"column:task_status"`
	Priority        string    `gorm:"column:priority"`
	DueDate         string    `gorm:"column:due_date"`
	DepartmentCode  string    `gorm:"column:department_code"`
	TeamsCode       string    `gorm:"column:teams_code"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedID       uint64    `gorm:"column:updated_id"`
	AssignedTo      uint      `gorm:"column:assigned_to"`
	AssignedAt      time.Time `gorm:"column:assigned_at"`
	StartedAt       time.Time `gorm:"column:started_at"`
	CompletedAt     time.Time `gorm:"column:completed_at"`
	SlaViolation    bool      `gorm:"column:sla_violation"`
	Comments        string    `gorm:"column:comments"`
	Attachments     string    `gorm:"column:attachments"`
	ReqTokenNo      string    `gorm:"column:req_token_no"`
	ReqRefType      string    `gorm:"column:req_ref_type"`
	ReqRefData      string    `gorm:"column:req_ref_data"`
	AssignedUser    string    `gorm:"column:assigned_user"`
}

// TableName overrides the default table name for GORM
func (TaskView) TableName() string {
	return "TaskView" // The name of your database view
}

// View Models Action View
type ActionView struct {
	ID                 uint      `gorm:"column:id"`
	CompanyID          uint      `gorm:"column:company_id"`
	TaskID             uint      `gorm:"column:task_id"`
	ActionName         string    `gorm:"column:action_name"`
	ActionDescription  string    `gorm:"column:action_description"`
	TranCode           string    `gorm:"column:tran_code"`
	SeqNo              uint      `gorm:"column:seq_no"`
	SlaDuration        uint      `gorm:"column:sla_duration"`
	SlaDurationType    string    `gorm:"column:sla_duration_type"`
	ActionStatus       string    `gorm:"column:action_status"`
	Priority           string    `gorm:"column:priority"`
	DueDate            string    `gorm:"column:due_date"`
	DepartmentCode     string    `gorm:"column:department_code"`
	TeamsCode          string    `gorm:"column:teams_code"`
	CreatedAt          time.Time `gorm:"column:created_at"`
	UpdatedID          uint64    `gorm:"column:updated_id"`
	TaskName           string    `gorm:"column:task_name"`
	TaskDescription    string    `gorm:"column:task_description"`
	ActionAssignmentId uint      `gorm:"column:action_assignment_id"`
	AssignedTo         uint      `gorm:"column:assigned_to"`
	AssignedAt         time.Time `gorm:"column:assigned_at"`
	StartedAt          time.Time `gorm:"column:started_at"`
	CompletedAt        time.Time `gorm:"column:completed_at"`
	SlaViolation       bool      `gorm:"column:sla_violation"`
	Comments           string    `gorm:"column:comments"`
	Attachments        string    `gorm:"column:attachments"`
	RequestId          uint      `gorm:"column:request_id"`
	ReqRefData         string    `gorm:"column:req_ref_data"`
	ReqRefType         string    `gorm:"column:req_ref_type"`
	ReqTokenNo         string    `gorm:"column:req_token_no"`
	AssignedUser       string    `gorm:"column:assigned_user"`
}

// TableName overrides the default table name for GORM
func (ActionView) TableName() string {
	return "ActionView" // The name of your database view
}

type WfComment struct {
	gorm.Model
	types.CModel
	TaskId         uint
	ActionId       uint
	AssignedUserID uint
	Sequence       uint
	Comment        string `gorm:"type:varchar(2500)"`
	CommentProof   string `gorm:"type:longtext"`
	// additional fields for record keeping
	//	AssignedTo     uint
	//	StatusSelected string `gorm:"type:varchar(2)"`
}
