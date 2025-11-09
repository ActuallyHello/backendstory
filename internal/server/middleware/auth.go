package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/services/auth"
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
				HandleError(w, r, errors.NewAccessError(nil, authMiddleware, "Missing authorization token!"))
				return
			}

			token := strings.TrimPrefix(authHeader, bearer)
			tokenUserInfo, err := authService.GetTokenUserInfo(ctx, token)
			roles := tokenUserInfo.Roles
			if err != nil {
				HandleError(w, r, errors.NewAccessError(err, authMiddleware, "Couldn't determinate roles for user!"))
				return
			}

			if !hasRequiredRole(roles, requiredRoles) {
				HandleError(w, r, errors.NewAccessError(nil, authMiddleware, "forbiden for user"))
			}

			next.ServeHTTP(w, r)
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
