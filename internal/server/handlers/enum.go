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
// @Tags enumerations
// @Accept json
// @Produce json
// @Param request body dto.EnumCreateRequest true "Данные для создания перечисления"
// @Success 201 {object} dto.EnumDTO "Созданное перечисление"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
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
	json.NewEncoder(w).Encode(h.toEnumDTO(enum))
}

// GetAll возвращает перечисления
// @Summary Получить перечисления
// @Description Возвращает перечисление по его идентификатору
// @Tags enumerations
// @Produce json
// @Success 200 {object} dto.EnumDTO "Найденные перечисление"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат"
// @Failure 404 {object} dto.ErrorResponse "Перечисление не найдено"
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
		dtos = append(dtos, h.toEnumDTO(enum))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetById возвращает перечисление по ID
// @Summary Получить перечисление по ID
// @Description Возвращает перечисление по его идентификатору
// @Tags enumerations
// @Produce json
// @Param id path int true "ID перечисления"
// @Success 200 {object} dto.EnumDTO "Найденное перечисление"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат ID"
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

	enum, err := h.enumerationService.GetById(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.toEnumDTO(enum))
}

// GetByCode возвращает перечисление по коду
// @Summary Получить перечисление по коду
// @Description Возвращает перечисление по его уникальному коду
// @Tags enumerations
// @Produce json
// @Param code path string true "Код перечисления"
// @Success 200 {object} dto.EnumDTO "Найденное перечисление"
// @Failure 400 {object} dto.ErrorResponse "Код не может быть пустым"
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
	json.NewEncoder(w).Encode(h.toEnumDTO(enum))
}

// // Update обновляет перечисление
// // @Summary Обновить перечисление
// // @Description Обновляет существующее перечисление
// // @Tags enumerations
// // @Accept json
// // @Produce json
// // @Param request body dto.EnumUpdateRequest true "Данные для обновления"
// // @Success 204 "Перечисление успешно обновлено"
// // @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// // @Failure 404 {object} dto.ErrorResponse "Перечисление не найдено"
// // @Failure 409 {object} dto.ErrorResponse "Конфликт данных (код уже существует)"
// // @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /enumerations [put]
// func (h *EnumHandler) Update(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	var req dto.EnumUpdateRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		middleware.HandleError(w, r, appErr.NewTechnicalError(err, enumHandlerCode, err.Error()))
// 		return
// 	}
// 	if err := h.validate.Struct(req); err != nil {
// 		details := common.CollectValidationDetails(err)
// 		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, enumHandlerCode, err.Error()), details)
// 		return
// 	}

// 	enum := entities.Enum{
// 		ID:    req.ID,
// 		Code:  req.Code,
// 		Label: req.Label,
// 	}
// 	_, err := h.enumerationService.Update(ctx, enum)
// 	if err != nil {
// 		middleware.HandleError(w, r, err)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// Delete удаляет перечисление
// @Summary Удалить перечисление
// @Description Удаляет перечисление (мягкое или полное удаление)
// @Tags enumerations
// @Produce json
// @Param id path int true "ID перечисления"
// @Param soft query bool false "Мягкое удаление (true/false)" default(true)
// @Success 204 "Перечисление успешно удалено"
// @Failure 400 {object} dto.ErrorResponse "Неверный формат ID или параметра soft"
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

	enum, err := h.enumerationService.GetById(ctx, uint(id))
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

func (h *EnumHandler) toEnumDTO(enum entities.Enum) dto.EnumDTO {
	return dto.EnumDTO{
		ID:    enum.ID,
		Code:  enum.Code,
		Label: enum.Label,
	}
}
