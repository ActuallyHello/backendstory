package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/ActuallyHello/backendstory/internal/services/auth"
)

type TokenCtx string
type UserInfoCtx string

const (
	authMiddleware = "AUTH_MIDDLEWARE_CODE"
	authorization  = "Authorization"
	bearer         = "Bearer "

	TokenCtxKey    TokenCtx    = "token"
	UserInfoCtxKey UserInfoCtx = "userInfo"
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

			ctx = context.WithValue(ctx, TokenCtxKey, token)
			ctx = context.WithValue(ctx, UserInfoCtxKey, tokenUserInfo)

			if !hasRequiredRole(roles, requiredRoles) {
				HandleError(w, r, errors.NewAccessError(nil, authMiddleware, "forbiden for user"))
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

func GetTokenCtx(ctx context.Context) (string, error) {
	tokenCtxKey, ok := ctx.Value(TokenCtxKey).(string)
	if !ok {
		return "", errors.NewLogicalError(nil, authMiddleware, "Couldn't find token in context")
	}
	return tokenCtxKey, nil
}

func GetUserInfoCtx(ctx context.Context) (dto.TokenUserInfo, error) {
	userInfoCtxKey, ok := ctx.Value(UserInfoCtxKey).(dto.TokenUserInfo)
	if !ok {
		return dto.TokenUserInfo{}, errors.NewLogicalError(nil, authMiddleware, "Couldn't find userInfo in context")
	}
	return userInfoCtxKey, nil
}
