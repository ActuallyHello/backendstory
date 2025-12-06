package person

import (
	"time"
)

// CreatePersonRequest represents request for creating person
// @Name CreatePersonRequest
type CreatePersonRequest struct {
	FirstName string `json:"firstname" validate:"required,min=2,max=50"`
	LastName  string `json:"lastname" validate:"required,min=2,max=50"`
	Phone     string `json:"phone" validate:"required"`
	UserLogin string `json:"user_login" validate:"required"`
}

// UpdatePersonRequest represents request for updating person
// @Name UpdatePersonRequest
type UpdatePersonRequest struct {
	ID        uint   `json:"id" validate:"required"`
	FirstName string `json:"firstName" validate:"omitempty,min=1,max=100"`
	LastName  string `json:"lastName" validate:"omitempty,min=1,max=100"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=20"`
	UserLogin string `json:"user_login" validate:"required"`
}

// PersonDTO represents person data transfer object
// @Name PersonDTO
type PersonDTO struct {
	ID        uint       `json:"id" validate:"required"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	UserLogin string `json:"user_login" validate:"required"`
}

func ToPersonDTO(person Person) PersonDTO {
	var deletedAt *time.Time
	if person.DeletedAt.Valid {
		deletedAt = &person.DeletedAt.Time
	}
	return PersonDTO{
		ID:        person.ID,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
		DeletedAt: deletedAt,

		Firstname: person.Firstname,
		Lastname:  person.Lastname,
		Phone:     person.Phone,
		UserLogin: person.UserLogin,
	}
}
