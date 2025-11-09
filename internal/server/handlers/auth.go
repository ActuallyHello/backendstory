package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
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

// Register регистрирует нового пользователя
// @Summary Регистрация пользователя
// @Description Создает нового пользователя и возвращает токен
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterUserRequest true "Данные для регистрации"
// @Success 200 {object} dto.JWT "JWT токен"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 409 {object} dto.ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /register [post]
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
		userDTO, getUserErr := h.authService.GetUserByUsername(ctx, req.Username)
		if getUserErr != nil {
			if !errors.Is(getUserErr, &appErr.LogicalError{}) {
				slog.Error("Couldn't compensate register user action!", "error", err)
				middleware.HandleError(w, r, getUserErr)
				return
			}
		}
		if userDTO.Username != "" {
			if deleteErr := h.authService.DeleteUser(ctx, userDTO.Username); deleteErr != nil {
				slog.Error("Couldn't compensate register user action! Deleted failed!", "error", err)
				middleware.HandleError(w, r, deleteErr)
				return
			}
		}
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(token)
}

// Login выполняет аутентификацию пользователя
// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя и возвращает токен
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} dto.JWT "JWT токен"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Неверные учетные данные"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /login [post]
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

// GetRoles возвращает список всех ролей
// @Summary Получить все роли
// @Description Возвращает список всех доступных ролей в системе
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} string "Список ролей"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/roles [get]
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

// GetUserRoles возвращает роли конкретного пользователя
// @Summary Получить роли пользователя
// @Description Возвращает список ролей для указанного пользователя
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Имя пользователя"
// @Success 200 {array} string "Список ролей пользователя"
// @Failure 400 {object} dto.ErrorResponse "Неверное имя пользователя"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/users/{username}/roles [get]
func (h *AuthHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := r.PathValue("username")
	if username == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, authHandlerCode, "username parameter missing"))
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

// GetUsers возвращает список всех пользователей
// @Summary Получить всех пользователей
// @Description Возвращает список всех зарегистрированных пользователей
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.UserDTO "Список пользователей"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/users [get]
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
