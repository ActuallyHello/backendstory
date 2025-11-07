package entities

type EnumValue struct {
	Base
	Code   string `gorm:"column:CODE"`
	Label  string `gorm:"column:LABEL"`
	EnumID uint   `gorm:"column:ENUMERATIONID"`
}

func (EnumValue) TableName() string {
	return "ENUMERATIONVALUE"
}
