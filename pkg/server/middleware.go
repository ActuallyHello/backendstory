package server

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/auth"
	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	authMiddleware = "AUTH_MIDDLEWARE_CODE"
	authorization  = "Authorization"
	bearer         = "Bearer "
)

func AuthMiddleware(authService auth.AuthService, requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authHeader := r.Header.Get(authorization)
			if authHeader == "" {
				core.HandleError(w, r, core.NewAccessError(nil, authMiddleware, "Missing authorization token!"))
				return
			}

			token := strings.TrimPrefix(authHeader, bearer)
			tokenUserInfo, err := authService.GetTokenUserInfo(ctx, token)
			roles := tokenUserInfo.Roles
			if err != nil {
				core.HandleError(w, r, core.NewAccessError(err, authMiddleware, "Couldn't determinate roles for user!"))
				return
			}

			ctx = context.WithValue(ctx, auth.TokenCtxKey, token)
			ctx = context.WithValue(ctx, auth.UserInfoCtxKey, tokenUserInfo)

			if !hasRequiredRole(roles, requiredRoles) {
				core.HandleError(w, r, core.NewAccessError(nil, authMiddleware, "forbiden for user"))
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func hasRequiredRole(userRoles, requiredRoles []string) bool {
	for _, need := range requiredRoles {
		if slices.Contains(userRoles, need) {
			return true
		}
	}
	return false
}
