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

// Register регистрирует нового пользователя
// @Summary Регистрация пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterUserRequest true "Данные для регистрации"
// @Success 201 {object} dto.UserDTO "Зарегистрированный пользователь"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 409 {object} dto.ErrorResponse "Пользователь с таким email уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/register [post]
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

	err := h.authService.RegisterUser(ctx, req.Email, req.Email, req.Password)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

// Login выполняет аутентификацию пользователя
// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя в систему и возвращает данные пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} dto.UserDTO "Аутентифицированный пользователь"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Неверные учетные данные"
// @Failure 404 {object} dto.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/login [post]
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

	tokenInfo, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewAccessError(err, authHandlerCode, err.Error()))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(tokenInfo)
}
