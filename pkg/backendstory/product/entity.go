package product

import (
	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/shopspring/decimal"
)

type Product struct {
	core.Base

	Label      string          `gorm:"column:LABEL"`
	Code       string          `gorm:"column:CODE"`
	Sku        string          `gorm:"column:SKU"`
	Price      decimal.Decimal `gorm:"column:PRICE"`
	Quantity   uint            `gorm:"column:QUANTITY"`
	CategoryID uint            `gorm:"column:CATEGORYID"`
	StatusID   uint            `gorm:"column:STATUSID"`
}

func (Product) TableName() string {
	return "PRODUCT"
}

func (Product) LocalTableName() string {
	return "Продукт"
}
