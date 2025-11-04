package entities

import (
	"database/sql"
	"time"
)

type Person struct {
	ID        uint         `gorm:"primaryKey;column:ID"`
	CreatedAt time.Time    `gorm:"column:CREATEDAT"`
	UpdatedAt time.Time    `gorm:"column:UPDATEDAT"`
	DeletedAt sql.NullTime `gorm:"column:DELETEDAT"`

	Firstname string `gorm:"column:FIRSTNAME"`
	Lastname  string `gorm:"column:LASTNAME"`
	Phone     string `gorm:"column:PHONE"`
	UserLogin string `gorm:"column:USERLOGIN"`
}

func (Person) TableName() string {
	return "PERSON"
}
