package entities

type CartItem struct {
	Base

	CartID    uint `gorm:"column:CARTID"`
	ProductID uint `gorm:"column:PRODUCTID"`
	StatusID  uint `gorm:"column:STATUSID"`
}

func (CartItem) TableName() string {
	return "CARTITEM"
}
