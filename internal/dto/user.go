package dto

import "time"

type RegisterUserRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=3,max=50"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserDTO struct {
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
}
