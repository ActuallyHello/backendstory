package cartitem

import "github.com/ActuallyHello/backendstory/pkg/core"

type CartItem struct {
	core.Base

	Quantity  uint `gorm:"column:QUANTITY"`
	CartID    uint `gorm:"column:CARTID"`
	ProductID uint `gorm:"column:PRODUCTID"`
	StatusID  uint `gorm:"column:STATUSID"`
}

func (CartItem) TableName() string {
	return "CARTITEM"
}

func (CartItem) LocalTableName() string {
	return "Элемент корзины"
}
