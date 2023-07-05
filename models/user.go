package models

import (
	"time"

	"github.com/shijith.chand/go-jwt/types"
)

type User struct {
	Id                      uint64        `gorm:"type:bigint;primaryKey;autoIncrement:true;"`
	Email                   string        `gorm:"type:varchar(50);unique"`
	Is_valid                types.BitBool `gorm:"type:bit(1)"`
	Name                    string        `gorm:"type:varchar(50)"`
	Password                string        `gorm:"type:varchar(100)"`
	Phone                   string        `gorm:"type:varchar(50);unique"`
	Auth_secret             string        `gorm:"type:varchar(50)"`
	Last_logged_inipaddress string        `gorm:"type:varchar(25)"`
	Last_logged_in_datime   time.Time
	DateFrom                string `gorm:"type:varchar(08)"`
	DateTo                  string `gorm:"type:varchar(08)"`
	Permissions             []Permission
	Profile                 string `gorm:"type:longtext"`
	VerficationCode         string `gorm:"type:varchar(10)"`
	LanguageID              uint
	Gender                  string `gorm:"type:varchar(1)"`
	UserGroupID             uint
	CompanyID               uint
	UserStatusID            uint
	BusinessDates           []BusinessDate
}
