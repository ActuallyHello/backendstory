package order

import (
	"database/sql"

	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	OrderStatus = "OrderStatus"

	PendingOrderStatus   = "InProgress"
	ApprovedOrderStatus  = "Approved"
	CancelledOrderStatus = "Cancelled"
)

type Order struct {
	core.Base

	Details   string        `gorm:"column:DETAILS"`
	StatusID  uint          `gorm:"column:STATUSID"`
	ClientID  uint          `gorm:"column:CLIENTID"`
	ManagerID sql.NullInt32 `gorm:"column:MANAGERID"`
}

func (Order) TableName() string {
	return "ORDERS"
}

func (Order) LocalTableName() string {
	return "Заказ"
}
