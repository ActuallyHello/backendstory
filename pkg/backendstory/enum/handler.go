package enum

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
)

const (
	enumHandlerCode = "ENUMERATION_HANDLER"
)

type EnumHandler struct {
	validate           *validator.Validate
	enumerationService EnumService
}

func NewEnumHandler(
	enumerationService EnumService,
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
// @Param request body EnumCreateRequest true "Данные для создания перечисления"
// @Success 201 {object} EnumDTO "Созданное перечисление"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} core.ErrorResponse "Перечисление с таким кодом уже существует"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations [post]
// @OperationId createEnumValue
func (h *EnumHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req EnumCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, enumHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, enumHandlerCode, err.Error()), details)
		return
	}

	enum := Enum{
		Code:  req.Code,
		Label: req.Label,
	}
	enum, err := h.enumerationService.Create(ctx, enum)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToEnumDTO(enum))
}

// GetAll возвращает все перечисления
// @Summary Получить все перечисления
// @Description Возвращает список всех перечислений в системе
// @Tags Enumerations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} EnumDTO "Список перечислений"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations [get]
// @OperationId getEnumAll
func (h *EnumHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	enums, err := h.enumerationService.GetAll(ctx)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]EnumDTO, 0, len(enums))
	for _, enum := range enums {
		dtos = append(dtos, ToEnumDTO(enum))
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
// @Success 200 {object} EnumDTO "Перечисление"
// @Failure 400 {object} core.ErrorResponse "Неверный ID"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Перечисление не найдено"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations/{id} [get]
// @OperationId getEnumById
func (h *EnumHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, enumHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, enumHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	enum, err := h.enumerationService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToEnumDTO(enum))
}

// GetByCode возвращает перечисление по коду
// @Summary Получить перечисление по коду
// @Description Возвращает перечисление по указанному коду
// @Tags Enumerations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Код перечисления"
// @Success 200 {object} EnumDTO "Перечисление"
// @Failure 400 {object} core.ErrorResponse "Неверный код"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Перечисление не найдено"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations/code/{code} [get]
// @OperationId getEnumByCode
func (h *EnumHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.PathValue("code")
	if code == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, enumHandlerCode, "Отсуствует параметр - код"))
		return
	}

	enum, err := h.enumerationService.GetByCode(ctx, code)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToEnumDTO(enum))
}

// GetWithSearchCriteria выполняет поиск перечислений по критериям
// @Summary Поиск перечислений
// @Description Выполняет поиск перечислений по заданным критериям
// @Tags Enumerations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.SearchCriteria true "Критерии поиска"
// @Success 200 {array} EnumDTO "Список найденных перечислений"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations/search [post]
// @OperationId searchEnum
func (h *EnumHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req core.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, enumHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, enumHandlerCode, err.Error()), details)
		return
	}

	enums, err := h.enumerationService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]EnumDTO, 0, len(enums))
	for _, enum := range enums {
		dtos = append(dtos, ToEnumDTO(enum))
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
// @Failure 400 {object} core.ErrorResponse "Неверный ID"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Перечисление не найдено"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /enumerations/{id} [delete]
// @OperationId deleteEnum
func (h *EnumHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, enumHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, enumHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	enum, err := h.enumerationService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}
	err = h.enumerationService.Delete(ctx, enum)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
