package entities

type Enum struct {
	Base
	Code  string `gorm:"column:CODE"`
	Label string `gorm:"column:LABEL"`
}

func (Enum) TableName() string {
	return "ENUMERATION"
}
