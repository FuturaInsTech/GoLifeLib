package models

import "time"

type ContHeader struct {
	ContractId       string `gorm:"type:varchar(20);primaryKey;"`
	ProductType      string `gorm:"type:varchar(10)"`
	StartDate        string `gorm:"type:varchar(8)"`
	EndDate          string `gorm:"type:varchar(8)"`
	Term             string `gorm:"type:varchar(3)"`
	Status           string `gorm:"type:varchar(10)"`
	BillingFrequency string `gorm:"type:varchar(2)"`
	PaymentMethod    string `gorm:"type:varchar(2)"`
	AgenyId          string `gorm:"type:varchar(20)"`
	IssueDate        string `gorm:"type:varchar(8)"`
	OwnerId          string `gorm:"type:varchar(20)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Converages       []Coverage
	CompanyID        uint
}
