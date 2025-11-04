package entities

import (
	"time"
)

type Enum struct {
	ID        uint      `gorm:"primaryKey;column:ID"`
	CreatedAt time.Time `gorm:"column:CREATEDAT"`
	UpdatedAt time.Time `gorm:"column:UPDATEDAT"`

	Code  string `gorm:"column:CODE"`
	Label string `gorm:"column:LABEL"`
}

func (Enum) TableName() string {
	return "ENUMERATION"
}
