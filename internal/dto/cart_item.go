package dto

import (
	"time"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
)

type CartItemCreateRequest struct {
	ProductID uint `json:"product_id" validate:"gte=1"`
	CartID    uint `json:"cart_id" validate:"gte=1"`
}

type CartItemDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ProductID uint      `json:"product_id"`
	CartID    uint      `json:"cart_id"`
}

func ToCartItemDTO(cartItem entities.CartItem) CartItemDTO {
	return CartItemDTO{
		ID:        cartItem.ID,
		CreatedAt: cartItem.CreatedAt,
		UpdatedAt: cartItem.UpdatedAt,
		ProductID: cartItem.ProductID,
		CartID:    cartItem.CartID,
	}
}
