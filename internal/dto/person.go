package dto

import "time"

type CreatePersonRequest struct {
	FirstName string `json:"firstname" validate:"required,min=2,max=50"`
	LastName  string `json:"lastname" validate:"required,min=2,max=50"`
	Phone     string `json:"phone" validate:"required"`
	UserLogin string `json:"user_login" validate:"required"`
}

type UpdatePersonRequest struct {
	ID        uint   `json:"id" validate:"required"`
	FirstName string `json:"firstName" validate:"omitempty,min=1,max=100"`
	LastName  string `json:"lastName" validate:"omitempty,min=1,max=100"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=20"`
	UserLogin string `json:"user_login" validate:"required"`
}

type PersonDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Phone     string `json:"phone"`
	UserLogin string `json:"user_login"`
}
