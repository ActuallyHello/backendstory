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
	"github.com/shopspring/decimal"
)

const (
	productHandlerCode = "PRODUCT_HANDLER"
)

type ProductHandler struct {
	validate              *validator.Validate
	producterationService services.ProductService
}

func NewProductHandler(
	producterationService services.ProductService,
) *ProductHandler {
	return &ProductHandler{
		validate:              validator.New(),
		producterationService: producterationService,
	}
}

// Create создает новое категории
// @Summary Создать категорию
// @Description Создает новое категорию в системе
// @Tags Producterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ProductCreateRequest true "Данные для создания категории"
// @Success 201 {object} dto.ProductDTO "Созданное категорию"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} dto.ErrorResponse "Продукт с таким кодом уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /productes [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.ProductCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, productHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, productHandlerCode, err.Error()), details)
		return
	}

	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, productHandlerCode, err.Error()))
		return
	}

	product := entities.Product{
		Code:       req.Code,
		Label:      req.Label,
		Sku:        req.Sku,
		Price:      price,
		Quantity:   req.Quantity,
		CategoryID: req.CategoryID,
		StatusID:   req.StatusID,
	}
	product, err = h.producterationService.Create(ctx, product)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToProductDTO(product))
}

// GetAll возвращает все категории
// @Summary Получить все категории
// @Description Возвращает список всех категория в системе
// @Tags Producterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ProductDTO "Список категория"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /productes [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.producterationService.GetAll(ctx)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.ProductDTO, 0, len(products))
	for _, product := range products {
		dtos = append(dtos, dto.ToProductDTO(product))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetById возвращает категорию по ID
// @Summary Получить категорию по ID
// @Description Возвращает категорию по указанному идентификатору
// @Tags Producterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID категории"
// @Success 200 {object} dto.ProductDTO "Продукт"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Продукт не найдено"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /productes/{id} [get]
func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, productHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, productHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	product, err := h.producterationService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ToProductDTO(product))
}

// GetByCode возвращает категорию по коду
// @Summary Получить категорию по коду
// @Description Возвращает категорию по указанному коду
// @Tags Producterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Код категории"
// @Success 200 {object} dto.ProductDTO "Продукт"
// @Failure 400 {object} dto.ErrorResponse "Неверный код"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Продукт не найдено"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /productes/code/{code} [get]
func (h *ProductHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.PathValue("code")
	if code == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, productHandlerCode, "ID parameter missing"))
		return
	}

	product, err := h.producterationService.GetByCode(ctx, code)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ToProductDTO(product))
}

// GetWithSearchCriteria выполняет поиск категория по критериям
// @Summary Поиск категория
// @Description Выполняет поиск категория по заданным критериям
// @Tags Producterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SearchCriteria true "Критерии поиска"
// @Success 200 {array} dto.ProductDTO "Список найденных категория"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /productes/search [post]
func (h *ProductHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, productHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, productHandlerCode, err.Error()), details)
		return
	}

	products, err := h.producterationService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.ProductDTO, 0, len(products))
	for _, product := range products {
		dtos = append(dtos, dto.ToProductDTO(product))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// Delete удаляет категорию
// @Summary Удалить категорию
// @Description Удаляет категорию по указанному идентификатору
// @Tags Producterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID категории"
// @Success 204 "Успешно удалено"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Продукт не найдено"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /productes/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, productHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, productHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	product, err := h.producterationService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}
	err = h.producterationService.Delete(ctx, product)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetByCode возвращает категорию по коду
// @Summary Получить категорию по коду
// @Description Возвращает категорию по указанному коду
// @Tags Producterations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "ID категории"
// @Success 200 {array} dto.ProductDTO "Продукты"
// @Failure 400 {object} dto.ErrorResponse "Неверный код"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Продукт не найдено"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /productes/category/{category_id} [get]
func (h *ProductHandler) GetByCategoryID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqCategoryID := r.PathValue("category_id")
	if reqCategoryID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, productHandlerCode, "ID parameter missing"))
		return
	}
	categoryID, err := strconv.Atoi(reqCategoryID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, enumValueHandlerCode, "CategoryID parameter must be integer!"+err.Error()))
		return
	}

	products, err := h.producterationService.GetByCategoryID(ctx, uint(categoryID))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.ProductDTO, 0, len(products))
	for _, product := range products {
		dtos = append(dtos, dto.ToProductDTO(product))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}
