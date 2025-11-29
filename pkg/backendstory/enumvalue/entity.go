package enumvalue

import "github.com/ActuallyHello/backendstory/pkg/core"

type EnumValue struct {
	core.Base
	Code   string `gorm:"column:CODE"`
	Label  string `gorm:"column:LABEL"`
	EnumID uint   `gorm:"column:ENUMERATIONID"`
}

func (EnumValue) TableName() string {
	return "ENUMERATIONVALUE"
}

func (EnumValue) LocalTableName() string {
	return "Значение перечисления"
}
