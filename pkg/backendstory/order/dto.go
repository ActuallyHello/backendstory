package order

import (
	"time"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
)

// OrderCreateRequest запрос на создание заказа
// @Name OrderCreateRequest
type OrderCreateRequest struct {
	ClientID    uint   `json:"client_id" validate:"required,min=1"`
	CartItemIDs []uint `json:"cart_item_ids" validate:"required"`
}

// OrderUpdateRequest запрос на создание заказа
// @Name OrderUpdateRequest
type OrderUpdateRequest struct {
	ID      uint   `json:"id" validate:"required,min=1"`
	Details string `json:"details" validate:"required"`
}

// OrderDTO представление заказа
// @Name OrderDTO
type OrderDTO struct {
	ID        uint                   `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Details   string                 `json:"details"`
	StatusDTO enumvalue.EnumValueDTO `json:"status_dto"`
	ClientID  uint                   `json:"client_id"`
	ManagerID *uint                  `json:"manager_id"`
}

func ToOrderDTO(order Order, status enumvalue.EnumValueDTO) OrderDTO {
	var managerID *uint
	if order.ManagerID.Valid {
		tempManagerID := uint(order.ManagerID.Int32)
		managerID = &tempManagerID
	}

	return OrderDTO{
		ID:        order.ID,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		Details:   order.Details,
		StatusDTO: status,
		ClientID:  order.ClientID,
		ManagerID: managerID,
	}
}
