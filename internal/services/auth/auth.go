package auth

import "context"

type AuthService interface {
	RegisterUser(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, username, password string) (map[string]any, error)
	GetRolesFromToken(ctx context.Context, token string) ([]string, error)
}
