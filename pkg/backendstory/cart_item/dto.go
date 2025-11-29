package cartitem

import (
	"time"
)

// CartItemCreateRequest represents request for creating cart item
// @Name CartItemCreateRequest
type CartItemCreateRequest struct {
	ProductID uint `json:"product_id" validate:"required,min=1"`
	CartID    uint `json:"cart_id" validate:"required,min=1"`
	Quantity  uint `json:"quantity" validate:"required,min=1"`
}

// CartItemUpdateRequest represents request for creating cart item
// @Name CartItemUpdateRequest
type CartItemUpdateRequest struct {
	CartItemID uint `json:"cart_item_id" validate:"required,min=1"`
	Quantity   uint `json:"quantity" validate:"required,min=1"`
}

// CartItemDTO represents cart item data transfer object
// @Name CartItemDTO
type CartItemDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ProductID uint      `json:"product_id"`
	CartID    uint      `json:"cart_id"`
	Quantity  uint      `json:"quantity"`
}

func ToCartItemDTO(cartItem CartItem) CartItemDTO {
	return CartItemDTO{
		ID:        cartItem.ID,
		CreatedAt: cartItem.CreatedAt,
		UpdatedAt: cartItem.UpdatedAt,
		ProductID: cartItem.ProductID,
		CartID:    cartItem.CartID,
		Quantity:  cartItem.Quantity,
	}
}
