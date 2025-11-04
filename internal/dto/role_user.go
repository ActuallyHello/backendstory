package dto

import "time"

type RoleUserCreateRequest struct {
	RoleID uint `json:"roleID" validate:"required"`
	UserID uint `json:"userID" validate:"required"`
}

type RoleUserUpdateRequest struct {
	ID     uint `json:"id" validate:"required"`
	RoleID uint `json:"roleID" validate:"required"`
	UserID uint `json:"userID" validate:"required"`
}

type RoleUserDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	RoleID    uint      `json:"roleID"`
	UserID    uint      `json:"userID"`
}
