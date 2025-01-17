package models

import "time"

type TaskView struct {
	ID              uint      `gorm:"column:id"`
	RequestID       uint      `gorm:"column:request_id"`
	TaskName        string    `gorm:"column:task_name"`
	TaskDescription string    `gorm:"column:task_description"`
	SeqNo           uint      `gorm:"column:seq_no"`
	SlaDuration     uint      `gorm:"column:sla_duration"`
	TaskStatus      string    `gorm:"column:task_status"`
	Priority        string    `gorm:"column:priority"`
	DueDate         string    `gorm:"column:due_date"`
	DepartmentCode  string    `gorm:"column:department_code"`
	TeamsCode       string    `gorm:"column:teams_code"`
	AssignedTo      uint      `gorm:"column:assigned_to"`
	AssignedAt      time.Time `gorm:"column:assigned_at"`
	StartedAt       time.Time `gorm:"column:started_at"`
	CompletedAt     time.Time `gorm:"column:completed_at"`
	SlaViolation    bool      `gorm:"column:sla_violation"`
	Comments        string    `gorm:"column:comments"`
	Attachments     string    `gorm:"column:attachments"`
	AssignedUser    uint      `gorm:"column:assigned_user"`
}

// TableName overrides the default table name for GORM
func (TaskView) TableName() string {
	return "TaskView" // The name of your database view
}
