package entities

type CartItem struct {
	Base

	Quantity  uint `gorm:"column:QUANTITY"`
	CartID    uint `gorm:"column:CARTID"`
	ProductID uint `gorm:"column:PRODUCTID"`
	StatusID  uint `gorm:"column:STATUSID"`
}

func (CartItem) TableName() string {
	return "CARTITEM"
}
