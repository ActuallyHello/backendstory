package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
)

const (
	authHandlerCode = "AUTH_HANDLER"
	authorization   = "Authorization"
	bearer          = "Bearer "
)

type AuthHandler struct {
	validate    *validator.Validate
	authService AuthService
}

func NewAuthHandler(
	authService AuthService,
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
// @Param request body RegisterUserRequest true "Данные для регистрации"
// @Success 200 {object} LoginResponse "Данные для входа"
// @Failure 409 {object} core.ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /register [post]
// @Id registerUser
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewValidationError(err, authHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewValidationError(err, authHandlerCode, err.Error()), details)
		return
	}

	// check if user with same email exists
	if _, err := h.authService.GetUserByEmail(ctx, req.Email); err == nil {
		core.HandleError(w, r, core.NewLogicalError(nil, authHandlerCode, "Пользователь с таким email уже зарегистрирован!"))
		return
	}

	err := h.authService.RegisterUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		slog.Error("Could'nt register user! Try to compensate...", "email", req.Email, "error", err)
		if deleteErr := h.authService.DeleteUser(ctx, req.Email); deleteErr != nil {
			slog.Error("Couldn't compensate user!", "email", req.Email, "error", deleteErr)
		} else {
			slog.Info("User was removed!", "email", req.Email)
		}
		core.HandleError(w, r, err)
		return
	}

	token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}
	tokenUserInfo, err := h.authService.GetTokenUserInfo(ctx, token.AccessToken)
	if err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	tokenResponse := LoginResponse{
		Token:    token,
		Username: tokenUserInfo.Username,
		Roles:    tokenUserInfo.Roles,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse)
}

// Login выполняет аутентификацию пользователя
// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя и возвращает токен
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} LoginResponse "Данные для входа"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Неверные учетные данные"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /login [post]
// @Id loginUser
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewValidationError(err, authHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewValidationError(err, authHandlerCode, err.Error()), details)
		return
	}

	token, err := h.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		core.HandleError(w, r, core.NewAccessError(err, authHandlerCode, err.Error()))
		return
	}

	tokenUserInfo, err := h.authService.GetTokenUserInfo(ctx, token.AccessToken)
	if err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	tokenResponse := LoginResponse{
		Token:    token,
		Username: tokenUserInfo.Username,
		Roles:    tokenUserInfo.Roles,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse)
}

// GetRoles возвращает список всех ролей
// @Summary Получить все роли
// @Description Возвращает список всех доступных ролей в системе
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} string "Список ролей"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/roles [get]
// @Id getRoles
func (h *AuthHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roles, err := h.authService.GetRoles(ctx)
	if err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, authHandlerCode, err.Error()))
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
// @Failure 400 {object} core.ErrorResponse "Неверное имя пользователя"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/users/{username}/roles [get]
// @Id getUserRoles
func (h *AuthHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := r.PathValue("username")
	if username == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, authHandlerCode, "Логин пользователя отсуствует"))
		return
	}

	roles, err := h.authService.GetRolesByUser(ctx, username)
	if err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, authHandlerCode, err.Error()))
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
// @Success 200 {array} UserDTO "Список пользователей"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/users [get]
// @Id getUsers
func (h *AuthHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.authService.GetUsers(ctx)
	if err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, authHandlerCode, err.Error()))
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
// @Success 200 {object} UserDTO "Пользователь"
// @Failure 400 {object} core.ErrorResponse "Неверное имя пользователя"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/users/{username} [get]
// @Id getUser
func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := r.PathValue("username")
	if username == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, authHandlerCode, "Логин пользователя отсуствует"))
		return
	}

	userDTO, err := h.authService.GetUserByEmail(ctx, username)
	if err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, authHandlerCode, err.Error()))
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
// @Success 200 {object} TokenUserInfo "Информация о пользователе из токена"
// @Failure 400 {object} core.ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/token [get]
// @Id getHeaderTokenInfo
func (h *AuthHandler) GetHeaderTokenInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get(authorization)
	if authHeader == "" {
		core.HandleError(w, r, core.NewAccessError(nil, authHandlerCode, "Отсуствует токен авторизации!"))
		return
	}

	token := strings.TrimPrefix(authHeader, bearer)
	tokenUserInfo, err := h.authService.GetTokenUserInfo(ctx, token)
	if err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, authHandlerCode, err.Error()))
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
// @Param request body TokenRequest true "Токен для проверки"
// @Success 200 {object} TokenUserInfo "Информация о пользователе из токена"
// @Failure 400 {object} core.ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/token [post]
// @Id getBodyTokenInfo
func (h *AuthHandler) GetBodyTokenInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var tokenReq TokenRequest

	if err := json.NewDecoder(r.Body).Decode(&tokenReq); err != nil {
		core.HandleError(w, r, core.NewValidationError(err, authHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(tokenReq); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewValidationError(err, authHandlerCode, err.Error()), details)
		return
	}

	tokenUserInfo, err := h.authService.GetTokenUserInfo(ctx, tokenReq.Token)
	if err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, authHandlerCode, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenUserInfo)
}
