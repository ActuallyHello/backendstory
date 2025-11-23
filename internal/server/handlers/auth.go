package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	appErr "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/ActuallyHello/backendstory/internal/server/handlers/common"
	"github.com/ActuallyHello/backendstory/internal/server/middleware"
	"github.com/ActuallyHello/backendstory/internal/services/auth"
	"github.com/go-playground/validator/v10"
)

const (
	authHandlerCode = "AUTH_HANDLER"
	authorization   = "Authorization"
	bearer          = "Bearer "
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
// @Success 200 {object} dto.LoginResponse "Данные для входа"
// @Failure 409 {object} dto.ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /register [post]
// @OperationId registerUser
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

	// check if user with same email exists
	if _, err := h.authService.GetUserByEmail(ctx, req.Email); err == nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, authHandlerCode, "User with this email already exists!"))
		return
	}

	err := h.authService.RegisterUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		slog.Error("Could'nt register user! Try to compensate...", "email", req.Email, "error", err)
		if deleteErr := h.authService.DeleteUser(ctx, req.Email); deleteErr != nil {
			slog.Error("Couldn't compensate user!", "email", req.Email, "error", deleteErr)
		} else {
			slog.Info("User was compensated!", "email", req.Email)
		}
		middleware.HandleError(w, r, err)
		return
	}

	token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}
	tokenUserInfo, err := h.authService.GetTokenUserInfo(ctx, token.AccessToken)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	loginResponse := dto.LoginResponse{
		Token:    token,
		Username: tokenUserInfo.Username,
		Roles:    tokenUserInfo.Roles,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
}

// Login выполняет аутентификацию пользователя
// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя и возвращает токен
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} dto.LoginResponse "Данные для входа"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Неверные учетные данные"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /login [post]
// @OperationId loginUser
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

	tokenUserInfo, err := h.authService.GetTokenUserInfo(ctx, token.AccessToken)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	loginResponse := dto.LoginResponse{
		Token:    token,
		Username: tokenUserInfo.Username,
		Roles:    tokenUserInfo.Roles,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
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
// @OperationId getRoles
func (h *AuthHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roles, err := h.authService.GetRoles(ctx)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
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
// @OperationId getUserRoles
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

	w.WriteHeader(http.StatusOK)
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
// @OperationId getUsers
func (h *AuthHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.authService.GetUsers(ctx)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetUser возвращает конкретного пользователя
// @Summary Получить пользователя
// @Description Возвращает пользователя по username
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Имя пользователя"
// @Success 200 {object} dto.UserDTO "Пользователь"
// @Failure 400 {object} dto.ErrorResponse "Неверное имя пользователя"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/users/{username} [get]
// @OperationId getUser
func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := r.PathValue("username")
	if username == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, authHandlerCode, "username parameter missing"))
		return
	}

	userDTO, err := h.authService.GetUserByEmail(ctx, username)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userDTO)
}

// GetHeaderTokenInfo возвращает информацию из текущего токена
// @Summary Получить информацию из текущего токена
// @Description Возвращает информацию из текущего токена
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.TokenUserInfo "Информация о пользователе из токена"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/token [get]
// @OperationId getHeaderTokenInfo
func (h *AuthHandler) GetHeaderTokenInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get(authorization)
	if authHeader == "" {
		middleware.HandleError(w, r, appErr.NewAccessError(nil, authHandlerCode, "Missing authorization token!"))
		return
	}

	token := strings.TrimPrefix(authHeader, bearer)
	tokenUserInfo, err := h.authService.GetTokenUserInfo(ctx, token)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenUserInfo)
}

// GetBodyTokenInfo возвращает информацию из токена
// @Summary Получить информацию из токена
// @Description Возвращает информацию из токена
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TokenRequest true "Токен для проверки"
// @Success 200 {object} dto.TokenUserInfo "Информация о пользователе из токена"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/token [post]
// @OperationId getBodyTokenInfo
func (h *AuthHandler) GetBodyTokenInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var tokenReq dto.TokenRequest

	if err := json.NewDecoder(r.Body).Decode(&tokenReq); err != nil {
		middleware.HandleError(w, r, appErr.NewValidationError(err, authHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(tokenReq); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewValidationError(err, authHandlerCode, err.Error()), details)
		return
	}

	tokenUserInfo, err := h.authService.GetTokenUserInfo(ctx, tokenReq.Token)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenUserInfo)
}
