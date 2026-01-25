package product

import (
	"database/sql"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/shopspring/decimal"
)

const (
	ProductStatus            = "ProductStatus"
	AvailableProductStatus   = "Available"
	UnAvailableProductStatus = "Unavailable"
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
	DeletedAt  sql.NullTime    `gorm:"column:DELETEDAT"`
	IsVisible  bool            `gorm:"column:ISVISIBLE"`
}

func (Product) TableName() string {
	return "PRODUCT"
}

func (Product) LocalTableName() string {
	return "Продукт"
}
