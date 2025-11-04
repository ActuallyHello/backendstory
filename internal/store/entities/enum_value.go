package entities

import (
	"time"
)

type EnumValue struct {
	ID        uint      `gorm:"primaryKey;column:ID"`
	CreatedAt time.Time `gorm:"column:CREATEDAT"`
	UpdatedAt time.Time `gorm:"column:UPDATEDAT"`

	Code   string `gorm:"column:CODE"`
	Label  string `gorm:"column:LABEL"`
	EnumID uint   `gorm:"column:ENUMERATIONID"`
}

func (EnumValue) TableName() string {
	return "ENUMERATIONVALUE"
}
