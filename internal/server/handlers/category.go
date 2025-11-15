package handlers

import (
	"database/sql"
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
	categoryHandlerCode = "CATEGORY_HANDLER"
)

type CategoryHandler struct {
	validate               *validator.Validate
	categoryerationService services.CategoryService
}

func NewCategoryHandler(
	categoryerationService services.CategoryService,
) *CategoryHandler {
	return &CategoryHandler{
		validate:               validator.New(),
		categoryerationService: categoryerationService,
	}
}

// Create создает новую категорию
// @Summary Создать категорию
// @Description Создает новую категорию в системе
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CategoryCreateRequest true "Данные для создания категории"
// @Success 201 {object} dto.CategoryDTO "Созданная категория"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} dto.ErrorResponse "Категория с таким кодом уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CategoryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, categoryHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, categoryHandlerCode, err.Error()), details)
		return
	}

	var categoryID sql.NullInt32
	if req.CategoryID != nil {
		categoryID = sql.NullInt32{
			Int32: int32(*req.CategoryID),
			Valid: true,
		}
	}

	category := entities.Category{
		Code:       req.Code,
		Label:      req.Label,
		CategoryID: categoryID,
	}
	category, err := h.categoryerationService.Create(ctx, category)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToCategoryDTO(category))
}

// GetAll возвращает все категории
// @Summary Получить все категории
// @Description Возвращает список всех категорий в системе
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.CategoryDTO "Список категорий"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.categoryerationService.GetAll(ctx)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.CategoryDTO, 0, len(categories))
	for _, category := range categories {
		dtos = append(dtos, dto.ToCategoryDTO(category))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetById возвращает категорию по ID
// @Summary Получить категорию по ID
// @Description Возвращает категорию по указанному идентификатору
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID категории"
// @Success 200 {object} dto.CategoryDTO "Категория"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Категория не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, categoryHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, categoryHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	category, err := h.categoryerationService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ToCategoryDTO(category))
}

// GetByCode возвращает категорию по коду
// @Summary Получить категорию по коду
// @Description Возвращает категорию по указанному коду
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Код категории"
// @Success 200 {object} dto.CategoryDTO "Категория"
// @Failure 400 {object} dto.ErrorResponse "Неверный код"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Категория не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/categories/code/{code} [get]
func (h *CategoryHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.PathValue("code")
	if code == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, categoryHandlerCode, "ID parameter missing"))
		return
	}

	category, err := h.categoryerationService.GetByCode(ctx, code)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ToCategoryDTO(category))
}

// GetWithSearchCriteria выполняет поиск категорий по критериям
// @Summary Поиск категорий
// @Description Выполняет поиск категорий по заданным критериям
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SearchCriteria true "Критерии поиска"
// @Success 200 {array} dto.CategoryDTO "Список найденных категорий"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/categories/search [post]
func (h *CategoryHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, categoryHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, categoryHandlerCode, err.Error()), details)
		return
	}

	categories, err := h.categoryerationService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.CategoryDTO, 0, len(categories))
	for _, category := range categories {
		dtos = append(dtos, dto.ToCategoryDTO(category))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// Delete удаляет категорию
// @Summary Удалить категорию
// @Description Удаляет категорию по указанному идентификатору
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID категории"
// @Success 204 "Успешно удалено"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Категория не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, categoryHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, categoryHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	category, err := h.categoryerationService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}
	err = h.categoryerationService.Delete(ctx, category)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetByCategoryID возвращает категории по родителю
// @Summary Получить категории по родителю
// @Description Возвращает категории по указанному родителю
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id path int true "ID родителя категории"
// @Success 200 {array} dto.CategoryDTO "Список категорий"
// @Failure 400 {object} dto.ErrorResponse "Неверная категория"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Категория не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/categories/category/{category_id} [get]
func (h *CategoryHandler) GetByCategoryID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqCategoryID := r.PathValue("category_id")
	if reqCategoryID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, categoryHandlerCode, "ID parameter missing"))
		return
	}
	categoryID, err := strconv.Atoi(reqCategoryID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, categoryHandlerCode, "CategoryID parameter must be integer!"+err.Error()))
		return
	}

	categories, err := h.categoryerationService.GetByCategoryID(ctx, uint(categoryID))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.CategoryDTO, 0, len(categories))
	for _, category := range categories {
		dtos = append(dtos, dto.ToCategoryDTO(category))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}
