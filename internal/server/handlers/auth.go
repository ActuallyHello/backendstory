package handlers

import (
	"encoding/json"
	"net/http"

	appErr "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/ActuallyHello/backendstory/internal/server/handlers/common"
	"github.com/ActuallyHello/backendstory/internal/server/middleware"
	"github.com/ActuallyHello/backendstory/internal/services/auth"
	"github.com/go-playground/validator/v10"
)

const (
	authHandlerCode = "AUTH_HANDLER"
)

type AuthHandler struct {
	validate    *validator.Validate
	authService auth.AuthService
}

func NewAuthHandler(
	authService auth.AuthService,
) *AuthHandler {
	return &AuthHandler{
		validate:    validator.New(),
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewValidationError(err, authHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewValidationError(err, authHandlerCode, err.Error()), details)
		return
	}

	token, err := h.authService.RegisterUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(token)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewValidationError(err, authHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewValidationError(err, authHandlerCode, err.Error()), details)
		return
	}

	token, err := h.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewAccessError(err, authHandlerCode, err.Error()))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(token)
}

func (h *AuthHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roles, err := h.authService.GetRoles(ctx)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(roles)
}

func (h *AuthHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := r.PathValue("username")
	if username == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, enumHandlerCode, "username parameter missing"))
		return
	}

	roles, err := h.authService.GetRolesByUser(ctx, username)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(roles)
}

func (h *AuthHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.authService.GetUsers(ctx)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(users)
}
