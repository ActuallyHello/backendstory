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
	enumHandlerCode = "ENUMERATION_HANDLER"
)

type EnumHandler struct {
	validate           *validator.Validate
	enumerationService services.EnumService
}

func NewEnumHandler(
	enumerationService services.EnumService,
) *EnumHandler {
	return &EnumHandler{
		validate:           validator.New(),
		enumerationService: enumerationService,
	}
}

// Create создает новое перечисление
// @Summary Создать перечисление
// @Description Создает новое перечисление в системе
// @Tags Enumerations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.EnumCreateRequest true "Данные для создания перечисления"
// @Success 201 {object} dto.EnumDTO "Созданное перечисление"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} dto.ErrorResponse "Перечисление с таким кодом уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations [post]
func (h *EnumHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.EnumCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, enumHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, enumHandlerCode, err.Error()), details)
		return
	}

	enum := entities.Enum{
		Code:  req.Code,
		Label: req.Label,
	}
	enum, err := h.enumerationService.Create(ctx, enum)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToEnumDTO(enum))
}

// GetAll возвращает все перечисления
// @Summary Получить все перечисления
// @Description Возвращает список всех перечислений в системе
// @Tags Enumerations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.EnumDTO "Список перечислений"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations [get]
func (h *EnumHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	enums, err := h.enumerationService.GetAll(ctx)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.EnumDTO, 0, len(enums))
	for _, enum := range enums {
		dtos = append(dtos, dto.ToEnumDTO(enum))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetById возвращает перечисление по ID
// @Summary Получить перечисление по ID
// @Description Возвращает перечисление по указанному идентификатору
// @Tags Enumerations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID перечисления"
// @Success 200 {object} dto.EnumDTO "Перечисление"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Перечисление не найдено"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations/{id} [get]
func (h *EnumHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, enumHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, enumHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	enum, err := h.enumerationService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ToEnumDTO(enum))
}

// GetByCode возвращает перечисление по коду
// @Summary Получить перечисление по коду
// @Description Возвращает перечисление по указанному коду
// @Tags Enumerations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Код перечисления"
// @Success 200 {object} dto.EnumDTO "Перечисление"
// @Failure 400 {object} dto.ErrorResponse "Неверный код"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Перечисление не найдено"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations/code/{code} [get]
func (h *EnumHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.PathValue("code")
	if code == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, enumHandlerCode, "ID parameter missing"))
		return
	}

	enum, err := h.enumerationService.GetByCode(ctx, code)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ToEnumDTO(enum))
}

// GetWithSearchCriteria выполняет поиск перечислений по критериям
// @Summary Поиск перечислений
// @Description Выполняет поиск перечислений по заданным критериям
// @Tags Enumerations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SearchCriteria true "Критерии поиска"
// @Success 200 {array} dto.EnumDTO "Список найденных перечислений"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations/search [post]
func (h *EnumHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, enumHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, enumHandlerCode, err.Error()), details)
		return
	}

	enums, err := h.enumerationService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.EnumDTO, 0, len(enums))
	for _, enum := range enums {
		dtos = append(dtos, dto.ToEnumDTO(enum))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// Delete удаляет перечисление
// @Summary Удалить перечисление
// @Description Удаляет перечисление по указанному идентификатору
// @Tags Enumerations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID перечисления"
// @Success 204 "Успешно удалено"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Перечисление не найдено"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations/{id} [delete]
func (h *EnumHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, enumHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, enumHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	enum, err := h.enumerationService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}
	err = h.enumerationService.Delete(ctx, enum)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
