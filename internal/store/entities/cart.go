package entities

type Cart struct {
	Base

	PersonID uint `gorm:"column:PERSONID"`
}

func (Cart) TableName() string {
	return "CART"
}
