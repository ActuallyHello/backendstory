package category

import (
	"database/sql"

	"github.com/ActuallyHello/backendstory/pkg/core"
)

type Category struct {
	core.Base

	Label      string        `gorm:"column:LABEL"`
	Code       string        `gorm:"column:CODE"`
	CategoryID sql.NullInt32 `gorm:"column:CATEGORYID"`
}

func (Category) TableName() string {
	return "CATEGORY"
}

func (Category) LocalTableName() string {
	return "Категория"
}
