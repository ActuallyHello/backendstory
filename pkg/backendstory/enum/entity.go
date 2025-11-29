package enum

import "github.com/ActuallyHello/backendstory/pkg/core"

type Enum struct {
	core.Base
	Code  string `gorm:"column:CODE"`
	Label string `gorm:"column:LABEL"`
}

func (Enum) TableName() string {
	return "ENUMERATION"
}

func (Enum) LocalTableName() string {
	return "Перечисление"
}
