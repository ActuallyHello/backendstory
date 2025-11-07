package entities

import "time"

type BaseEntity interface {
	TableName() string
	GetID() uint
}

type Base struct {
	ID        uint      `gorm:"primaryKey;column:ID"`
	CreatedAt time.Time `gorm:"column:CREATEDAT"`
	UpdatedAt time.Time `gorm:"column:UPDATEDAT"`
}

func (b Base) GetID() uint {
	return b.ID
}
