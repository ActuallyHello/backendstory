package cart

import "github.com/ActuallyHello/backendstory/pkg/core"

type Cart struct {
	core.Base

	PersonID uint `gorm:"column:PERSONID"`
}

func (Cart) TableName() string {
	return "CART"
}

func (Cart) LocalTableName() string {
	return "Корзина"
}
