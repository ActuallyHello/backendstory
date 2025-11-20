package dto

import "time"

// RegisterUserRequest represents request for user registration
// @Name RegisterUserRequest
type RegisterUserRequest struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=3,max=50"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// LoginRequest represents request for user login
// @Name LoginRequest
type LoginRequest struct {
	Login    string `json:"login" validate:"required"` // could be username or email
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents response for user login
// @Name LoginResponse
type LoginResponse struct {
	Token    JWT      `json:"token"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

// UserDTO represents user data transfer object
// @Name UserDTO
type UserDTO struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// JWT represents JWT tokens response
// @Name JWT
type JWT struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
}

// TokenUserInfo represents user information from token
// @Name TokenUserInfo
type TokenUserInfo struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

// TokenRequest represents token verification request
// @Name TokenRequest
type TokenRequest struct {
	Token string `json:"token" validate:"required"`
}
