package person

import (
	"database/sql"

	"github.com/ActuallyHello/backendstory/pkg/core"
)

type Person struct {
	core.Base
	DeletedAt sql.NullTime `gorm:"column:DELETEDAT"`

	Firstname string `gorm:"column:FIRSTNAME"`
	Lastname  string `gorm:"column:LASTNAME"`
	Phone     string `gorm:"column:PHONE"`
	UserLogin string `gorm:"column:USERLOGIN"`
}

func (Person) TableName() string {
	return "PERSON"
}

func (Person) LocalTableName() string {
	return "Клиент"
}
