package entities

import "database/sql"

type Category struct {
	Base

	Label      string        `gorm:"column:LABEL"`
	Code       string        `gorm:"column:CODE"`
	CategoryID sql.NullInt32 `gorm:"column:CATEGORYID"`
}

func (Category) TableName() string {
	return "CATEGORY"
}
