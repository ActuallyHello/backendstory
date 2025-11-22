package dto

import (
	"time"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
)

type CartCreateRequest struct {
	PersonID uint `json:"person_id" validate:"required,min=1,max=50"`
}

type CartDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PersonID  uint      `json:"person_id"`
}

func ToCartDTO(cart entities.Cart) CartDTO {
	return CartDTO{
		ID:        cart.ID,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
		PersonID:  cart.PersonID,
	}
}
