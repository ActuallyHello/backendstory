package handlers

import (
	"encoding/json"
	"fmt"
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

func NewEnumValuenHandler(
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
	fmt.Println("SHOW ENUM VALUE", enumValue)
	enumValue, err := h.enumValueService.Create(ctx, enumValue)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h.toEnumValueDTO(enumValue))
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
func (h *EnumValueHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	enumValues, err := h.enumValueService.GetAll(ctx)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.toEnumValueDTOs(enumValues))
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
	json.NewEncoder(w).Encode(h.toEnumValueDTO(enumValue))
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.toEnumValueDTOs(enumValues))
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

func (h *EnumValueHandler) toEnumValueDTO(enumValue entities.EnumValue) dto.EnumValueDTO {
	return dto.EnumValueDTO{
		ID:     enumValue.ID,
		Code:   enumValue.Code,
		Label:  enumValue.Label,
		EnumID: enumValue.EnumID,
	}
}

func (h *EnumValueHandler) toEnumValueDTOs(values []entities.EnumValue) []dto.EnumValueDTO {
	dtos := make([]dto.EnumValueDTO, len(values))
	for i := 0; i < len(values); i++ {
		dtos[i] = h.toEnumValueDTO(values[i])
	}
	return dtos
}
