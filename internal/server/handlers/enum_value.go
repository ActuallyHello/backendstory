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
	enumValueHandlerCode = "ENUMERATION_VALUE_HANDLER"
)

type EnumValueHandler struct {
	validate         *validator.Validate
	enumValueService services.EnumValueService
}

func NewEnumValueHandler(
	enumValueService services.EnumValueService,
) *EnumValueHandler {
	return &EnumValueHandler{
		validate:         validator.New(),
		enumValueService: enumValueService,
	}
}

// Create создает новое значение перечисления
// @Summary Создать значение перечисления
// @Description Создает новое значение перечисления
// @Tags Enumeration Values
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.EnumValueCreateRequest true "Данные для создания значения перечисления"
// @Success 201 {object} dto.EnumValueDTO "Созданное значение перечисления"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumeration-values [post]
// @OperationId createEnumValue
func (h *EnumValueHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.EnumValueCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, enumValueHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, enumValueHandlerCode, err.Error()), details)
		return
	}

	enumValue := entities.EnumValue{
		Code:   req.Code,
		Label:  req.Label,
		EnumID: req.EnumID,
	}
	enumValue, err := h.enumValueService.Create(ctx, enumValue)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToEnumValueDTO(enumValue))
}

// GetAll возвращает все значения перечислений
// @Summary Получить все значения перечислений
// @Description Возвращает список всех значений перечислений
// @Tags Enumeration Values
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.EnumValueDTO "Список значений перечислений"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumeration-values [get]
// @OperationId getEnumValueAll
func (h *EnumValueHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	enumValues, err := h.enumValueService.GetAll(ctx)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.EnumValueDTO, 0, len(enumValues))
	for _, enumValue := range enumValues {
		dtos = append(dtos, dto.ToEnumValueDTO(enumValue))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetById возвращает значение перечисления по ID
// @Summary Получить значение перечисления по ID
// @Description Возвращает значение перечисления по указанному идентификатору
// @Tags Enumeration Values
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID значения перечисления"
// @Success 200 {object} dto.EnumValueDTO "Значение перечисления"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Значение перечисления не найдено"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumeration-values/{id} [get]
// @OperationId getEnumValueById
func (h *EnumValueHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, enumValueHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, enumValueHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	enumValue, err := h.enumValueService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ToEnumValueDTO(enumValue))
}

// GetByEnumId возвращает значения перечисления по ID перечисления
// @Summary Получить значения по ID перечисления
// @Description Возвращает все значения для указанного перечисления
// @Tags Enumeration Values
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param enumeration_id path int true "ID перечисления"
// @Success 200 {array} dto.EnumValueDTO "Список значений перечисления"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID перечисления"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumeration-values/enumeration/{enumeration_id} [get]
// @OperationId getEnumValueByEnumId
func (h *EnumValueHandler) GetByEnumId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqEnumID := r.PathValue("enumeration_id")
	if reqEnumID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, enumValueHandlerCode, "ID parameter missing"))
		return
	}
	enumID, err := strconv.Atoi(reqEnumID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, enumValueHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	enumValues, err := h.enumValueService.GetByEnumID(ctx, uint(enumID))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.EnumValueDTO, 0, len(enumValues))
	for _, enumValue := range enumValues {
		dtos = append(dtos, dto.ToEnumValueDTO(enumValue))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// Delete удаляет значение перечисления
// @Summary Удалить значение перечисления
// @Description Удаляет значение перечисления по указанному идентификатору
// @Tags Enumeration Values
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID значения перечисления"
// @Success 204 "Успешно удалено"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Значение перечисления не найдено"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumeration-values/{id} [delete]
// @OperationId deleteEnumValue
func (h *EnumValueHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, enumValueHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, enumValueHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	enumValue, err := h.enumValueService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}
	err = h.enumValueService.Delete(ctx, enumValue)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetWithSearchCriteria выполняет поиск значений перечислений по критериям
// @Summary Поиск значений перечислений
// @Description Выполняет поиск значений перечислений по заданным критериям
// @Tags Enumeration Values
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SearchCriteria true "Критерии поиска"
// @Success 200 {array} dto.EnumValueDTO "Список найденных значений перечислений"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumeration-values/search [post]
// @OperationId searchEnumValue
func (h *EnumValueHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, enumValueHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, enumValueHandlerCode, err.Error()), details)
		return
	}

	enumValues, err := h.enumValueService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.EnumValueDTO, 0, len(enumValues))
	for _, enumValue := range enumValues {
		dtos = append(dtos, dto.ToEnumValueDTO(enumValue))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}
