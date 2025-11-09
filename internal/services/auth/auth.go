package auth

import (
	"context"

	"github.com/ActuallyHello/backendstory/internal/dto"
)

type AuthService interface {
	RegisterUser(ctx context.Context, username, email, password string) (dto.JWT, error)
	DeleteUser(ctx context.Context, username string) error
	Login(ctx context.Context, username, password string) (dto.JWT, error)

	RefreshToken(ctx context.Context, refreshToken string) (dto.JWT, error)

	GetUserByUsername(ctx context.Context, username string) (dto.UserDTO, error)
	GetUsers(ctx context.Context) ([]dto.UserDTO, error)
	GetRoles(ctx context.Context) ([]string, error)
	GetRolesByUser(ctx context.Context, username string) ([]string, error)
	GetRolesFromToken(ctx context.Context, token string) ([]string, error)
}
