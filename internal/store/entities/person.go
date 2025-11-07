package entities

import (
	"database/sql"
)

type Person struct {
	Base
	DeletedAt sql.NullTime `gorm:"column:DELETEDAT"`

	Firstname string `gorm:"column:FIRSTNAME"`
	Lastname  string `gorm:"column:LASTNAME"`
	Phone     string `gorm:"column:PHONE"`
	UserLogin string `gorm:"column:USERLOGIN"`
}

func (Person) TableName() string {
	return "PERSON"
}
