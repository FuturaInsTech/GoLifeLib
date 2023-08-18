package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type BitBool bool
type ExtraData map[string]interface{}
type LockedType uint8

const (
	Contract LockedType = iota
	Client
)

func (bb BitBool) Value() (driver.Value, error) {
	return bool(bb), nil
}

func (bb *BitBool) Scan(src interface{}) error {
	if src == nil {
		// MySql NULL value turns into false
		*bb = false
		return nil
	}
	bs, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Not byte slice!")
	}
	*bb = bs[0] == 1
	return nil
}

type CModel struct {
	IsValid   string `gorm:"type:varchar(01)"`
	UpdatedID uint64
	Tranno    uint
	Session   string `gorm:"type:varchar(15)"`
	SessionID uint
	CompanyID uint
}

func (ed ExtraData) Value() (driver.Value, error) {

	data, err := json.Marshal(ed)
	return string(data), err

}

func (ed *ExtraData) Scan(src interface{}) error {

	bs, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("Not byte slice!")
	}

	return json.Unmarshal(bs, ed)
}

func (u *LockedType) Scan(value interface{}) error { *u = LockedType(value.(int64)); return nil }
func (u LockedType) Value() (driver.Value, error)  { return int64(u), nil }
