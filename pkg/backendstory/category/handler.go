package category

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
)

const (
	categoryHandlerCode = "CATEGORY_HANDLER"
)

type CategoryHandler struct {
	validate               *validator.Validate
	categoryerationService CategoryService
}

func NewCategoryHandler(
	categoryerationService CategoryService,
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
// @Param request body CategoryCreateRequest true "Данные для создания категории"
// @Success 201 {object} CategoryDTO "Созданная категория"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} core.ErrorResponse "Категория с таким кодом уже существует"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories [post]
// @ID createCategory
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CategoryCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, categoryHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, categoryHandlerCode, err.Error()), details)
		return
	}

	var categoryID sql.NullInt32
	if req.CategoryID != nil {
		categoryID = sql.NullInt32{
			Int32: int32(*req.CategoryID),
			Valid: true,
		}
	}

	category := Category{
		Code:       req.Code,
		Label:      req.Label,
		CategoryID: categoryID,
	}
	category, err := h.categoryerationService.Create(ctx, category)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToCategoryDTO(category))
}

// GetAll возвращает все категории
// @Summary Получить все категории
// @Description Возвращает список всех категорий в системе
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} CategoryDTO "Список категорий"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories [get]
// @ID getCategoryAll
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.categoryerationService.GetAll(ctx)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]CategoryDTO, 0, len(categories))
	for _, category := range categories {
		dtos = append(dtos, ToCategoryDTO(category))
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
// @Success 200 {object} CategoryDTO "Категория"
// @Failure 400 {object} core.ErrorResponse "Неверный ID"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Категория не найдена"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories/{id} [get]
// @ID getCategoryById
func (h *CategoryHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, categoryHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, categoryHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	category, err := h.categoryerationService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToCategoryDTO(category))
}

// GetByCode возвращает категорию по коду
// @Summary Получить категорию по коду
// @Description Возвращает категорию по указанному коду
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Код категории"
// @Success 200 {object} CategoryDTO "Категория"
// @Failure 400 {object} core.ErrorResponse "Неверный код"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Категория не найдена"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories/code/{code} [get]
// @ID getCategoryByCode
func (h *CategoryHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.PathValue("code")
	if code == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, categoryHandlerCode, "ID parameter missing"))
		return
	}

	category, err := h.categoryerationService.GetByCode(ctx, code)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToCategoryDTO(category))
}

// GetWithSearchCriteria выполняет поиск категорий по критериям
// @Summary Поиск категорий
// @Description Выполняет поиск категорий по заданным критериям
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.SearchCriteria true "Критерии поиска"
// @Success 200 {array} CategoryDTO "Список найденных категорий"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories/search [post]
// @ID getCategorySearch
func (h *CategoryHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req core.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, categoryHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, categoryHandlerCode, err.Error()), details)
		return
	}

	categories, err := h.categoryerationService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]CategoryDTO, 0, len(categories))
	for _, category := range categories {
		dtos = append(dtos, ToCategoryDTO(category))
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
// @Failure 400 {object} core.ErrorResponse "Неверный ID"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Категория не найдена"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories/{id} [delete]
// @ID deleteCategory
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, categoryHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, categoryHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	category, err := h.categoryerationService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}
	err = h.categoryerationService.Delete(ctx, category)
	if err != nil {
		core.HandleError(w, r, err)
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
// @Success 200 {array} CategoryDTO "Список категорий"
// @Failure 400 {object} core.ErrorResponse "Неверная категория"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Категория не найдена"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories/category/{category_id} [get]
// @ID getCategoryByCategoryId
func (h *CategoryHandler) GetByCategoryID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqCategoryID := r.PathValue("category_id")
	if reqCategoryID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, categoryHandlerCode, "ID parameter missing"))
		return
	}
	categoryID, err := strconv.Atoi(reqCategoryID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, categoryHandlerCode, "CategoryID parameter must be integer!"+err.Error()))
		return
	}

	categories, err := h.categoryerationService.GetByCategoryID(ctx, uint(categoryID))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]CategoryDTO, 0, len(categories))
	for _, category := range categories {
		dtos = append(dtos, ToCategoryDTO(category))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}
