package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	appErr "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/ActuallyHello/backendstory/internal/server/handlers/common"
	"github.com/ActuallyHello/backendstory/internal/server/middleware"
	"github.com/ActuallyHello/backendstory/internal/services"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/go-playground/validator/v10"
)

const (
	roleHandlerCode = "ROLE_HANDLER"
)

type RoleHandler struct {
	validate    *validator.Validate
	roleService services.RoleService
}

func NewRoleHandler(
	roleService services.RoleService,
) *RoleHandler {
	return &RoleHandler{
		validate:    validator.New(),
		roleService: roleService,
	}
}

// Create создает новую роль
// @Summary Создать роль
// @Description Создает новую роль в системе
// @Tags roles
// @Accept json
// @Produce json
// @Param request body dto.RoleCreateRequest true "Данные для создания роли"
// @Success 201 {object} dto.RoleDTO "Созданная роль"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 409 {object} dto.ErrorResponse "Роль с таким кодом уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles [post]
func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.RoleCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, roleHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, roleHandlerCode, err.Error()), details)
		return
	}

	role := entities.Role{
		Code:  req.Code,
		Label: req.Label,
	}
	role, err := h.roleService.Create(ctx, role)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h.toRoleDTO(role))
}

// GetById возвращает роли
// @Summary Получить роли
// @Description Возвращает роль по ее идентификатору
// @Tags roles
// @Produce json
// @Success 200 {object} dto.RoleDTO "Найденная роль"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат ID"
// @Failure 404 {object} dto.ErrorResponse "Роль не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles [get]
func (h *RoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roles, err := h.roleService.GetAll(ctx)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.RoleDTO, 0, len(roles))
	for _, role := range roles {
		dtos = append(dtos, h.toRoleDTO(role))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetById возвращает роль по ID
// @Summary Получить роль по ID
// @Description Возвращает роль по ее идентификатору
// @Tags roles
// @Produce json
// @Param id path int true "ID роли"
// @Success 200 {object} dto.RoleDTO "Найденная роль"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат ID"
// @Failure 404 {object} dto.ErrorResponse "Роль не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles/{id} [get]
func (h *RoleHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, roleHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, roleHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	role, err := h.roleService.GetById(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.toRoleDTO(role))
}

// GetByCode возвращает роль по коду
// @Summary Получить роль по коду
// @Description Возвращает роль по ее уникальному коду
// @Tags roles
// @Produce json
// @Param code path string true "Код роли"
// @Success 200 {object} dto.RoleDTO "Найденная роль"
// @Failure 400 {object} dto.ErrorResponse "Код не может быть пустым"
// @Failure 404 {object} dto.ErrorResponse "Роль не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles/code/{code} [get]
func (h *RoleHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.PathValue("code")
	if code == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, roleHandlerCode, "Code parameter missing"))
		return
	}

	role, err := h.roleService.GetByCode(ctx, code)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.toRoleDTO(role))
}

// Update обновляет роль
// @Summary Обновить роль
// @Description Обновляет существующую роль
// @Tags roles
// @Accept json
// @Produce json
// @Param request body dto.RoleUpdateRequest true "Данные для обновления"
// @Success 204 "Роль успешно обновлена"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 404 {object} dto.ErrorResponse "Роль не найдена"
// @Failure 409 {object} dto.ErrorResponse "Конфликт данных (код уже существует)"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles [put]
func (h *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.RoleUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, roleHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, roleHandlerCode, err.Error()), details)
		return
	}

	role := entities.Role{
		ID:    req.ID,
		Code:  req.Code,
		Label: req.Label,
	}
	_, err := h.roleService.Update(ctx, role)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete удаляет роль
// @Summary Удалить роль
// @Description Удаляет роль
// @Tags roles
// @Produce json
// @Param id path int true "ID роли"
// @Success 204 "Роль успешно удалена"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат ID"
// @Failure 404 {object} dto.ErrorResponse "Роль не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /roles/{id} [delete]
func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, roleHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, roleHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	role, err := h.roleService.GetById(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}
	err = h.roleService.Delete(ctx, role)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RoleHandler) toRoleDTO(role entities.Role) dto.RoleDTO {
	return dto.RoleDTO{
		ID:        role.ID,
		CreatedAt: role.CreatedAt,
		Code:      role.Code,
		Label:     role.Label,
	}
}
