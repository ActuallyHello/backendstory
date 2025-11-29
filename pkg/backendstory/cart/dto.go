package cart

import (
	"time"
)

// CartCreateRequest represents request for creating cart
// @Name CartCreateRequest
type CartCreateRequest struct {
	PersonID uint `json:"person_id" validate:"required,min=1"`
}

// CartDTO represents cart data transfer object
// @Name CartDTO
type CartDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PersonID  uint      `json:"person_id"`
}

func ToCartDTO(cart Cart) CartDTO {
	return CartDTO{
		ID:        cart.ID,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
		PersonID:  cart.PersonID,
	}
}
