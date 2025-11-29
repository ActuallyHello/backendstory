package productmedia

import "github.com/ActuallyHello/backendstory/pkg/core"

type ProductMedia struct {
	core.Base

	Link      string `gorm:"column:LINK"`
	ProductID uint   `gorm:"column:PRODUCTID"`
}

func (ProductMedia) TableName() string {
	return "PRODUCTMEDIA"
}

func (ProductMedia) LocalTableName() string {
	return "Картинка товара"
}
