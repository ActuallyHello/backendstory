package orderitem

import (
	"time"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
)

// OrderItemCreateRequest запрос на создание элемента заказа
// @Name OrderItemCreateRequest
type OrderItemCreateRequest struct {
	OrderID    uint `json:"order_id" validate:"required,min=1"`
	CartItemID uint `json:"cart_item_id" validate:"required,min=1"`
}

// OrderDTO представление заказа
// @Name OrderDTO
type OrderItemDTO struct {
	ID         uint                   `json:"id"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	StatusDTO  enumvalue.EnumValueDTO `json:"status_dto"`
	OrderID    uint                   `json:"order_id"`
	CartItemId uint                   `json:"cart_item_id"`
}

func ToOrderItemDTO(orderItem OrderItem, status enumvalue.EnumValueDTO) OrderItemDTO {
	return OrderItemDTO{
		ID:         orderItem.ID,
		CreatedAt:  orderItem.CreatedAt,
		UpdatedAt:  orderItem.UpdatedAt,
		StatusDTO:  status,
		OrderID:    orderItem.OrderID,
		CartItemId: orderItem.CartItemID,
	}
}
