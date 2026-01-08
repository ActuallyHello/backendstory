package orderitem

import (
	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	OrderItemStatus = "OrderItemStatus"

	PendingOrderItemStatus   = "InProgress"
	ApprovedOrderItemStatus  = "Approved"
	CancelledOrderItemStatus = "Cancelled"
)

type OrderItem struct {
	core.Base

	StatusID   uint `gorm:"column:STATUSID"`
	OrderID    uint `gorm:"column:ORDERID"`
	CartItemID uint `gorm:"column:CARTITEMID"`
}

func (OrderItem) TableName() string {
	return "ORDERITEM"
}

func (OrderItem) LocalTableName() string {
	return "Элемента заказа"
}
