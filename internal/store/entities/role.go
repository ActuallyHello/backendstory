package entities

import (
	"time"
)

type Role struct {
	ID        uint      `gorm:"primaryKey;column:ID"`
	CreatedAt time.Time `gorm:"column:CREATEDAT"`

	Code  string `gorm:"column:CODE"`
	Label string `gorm:"column:LABEL"`
}

func (Role) TableName() string {
	return "ROLE_"
}
