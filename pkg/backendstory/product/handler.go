package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/category"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

const (
	productHandlerCode = "PRODUCT_HANDLER"
)

type ProductHandler struct {
	validate              *validator.Validate
	producterationService ProductService
	enumValueService      enumvalue.EnumValueService
	categoryService       category.CategoryService
}

func NewProductHandler(
	producterationService ProductService,
	enumValueService enumvalue.EnumValueService,
	categoryService category.CategoryService,
) *ProductHandler {
	return &ProductHandler{
		validate:              validator.New(),
		producterationService: producterationService,
		enumValueService:      enumValueService,
		categoryService:       categoryService,
	}
}

// Create создает новый продукт
// @Summary Создать продукт
// @Description Создает новый продукт в системе
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ProductCreateRequest true "Данные для создания продукта"
// @Success 201 {object} ProductDTO "Созданный продукт"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} core.ErrorResponse "Продукт с таким кодом уже существует"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products [post]
// @OperationId createProduct
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ProductCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, productHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, productHandlerCode, err.Error()), details)
		return
	}

	// check price
	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, productHandlerCode, err.Error()))
		return
	}
	// check status
	if _, err := h.enumValueService.GetByID(ctx, req.StatusID); err != nil {
		core.HandleError(w, r, err)
		return
	}
	// check category
	if _, err := h.categoryService.GetByID(ctx, req.CategoryID); err != nil {
		core.HandleError(w, r, err)
		return
	}

	product := Product{
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
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToProductDTO(product))
}

// GetAll возвращает все продукты
// @Summary Получить все продукты
// @Description Возвращает список всех продуктов в системе
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ProductDTO "Список продуктов"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products [get]
// @OperationId getProductAll
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.producterationService.GetAll(ctx)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]ProductDTO, 0, len(products))
	for _, product := range products {
		dtos = append(dtos, ToProductDTO(product))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetById возвращает продукт по ID
// @Summary Получить продукт по ID
// @Description Возвращает продукт по указанному идентификатору
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID продукта"
// @Success 200 {object} ProductDTO "Продукт"
// @Failure 400 {object} core.ErrorResponse "Неверный ID"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Продукт не найден"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/{id} [get]
// @OperationId getProductById
func (h *ProductHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, productHandlerCode, "Отсутствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, productHandlerCode, "ИД параметр должен быть числовым! "+err.Error()))
		return
	}

	product, err := h.producterationService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToProductDTO(product))
}

// GetByCode возвращает продукт по коду
// @Summary Получить продукт по коду
// @Description Возвращает продукт по указанному коду
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Код продукта"
// @Success 200 {object} ProductDTO "Продукт"
// @Failure 400 {object} core.ErrorResponse "Неверный код"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Продукт не найден"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/code/{code} [get]
// @OperationId getProductByCode
func (h *ProductHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := r.PathValue("code")
	if code == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, productHandlerCode, "Код параметра обзятаельна"))
		return
	}

	product, err := h.producterationService.GetByCode(ctx, code)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToProductDTO(product))
}

// GetWithSearchCriteria выполняет поиск продуктов по критериям
// @Summary Поиск продуктов
// @Description Выполняет поиск продуктов по заданным критериям
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.SearchCriteria true "Критерии поиска"
// @Success 200 {array} ProductDTO "Список найденных продуктов"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/search [post]
// @OperationId searchProduct
func (h *ProductHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req core.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, productHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, productHandlerCode, err.Error()), details)
		return
	}

	products, err := h.producterationService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]ProductDTO, 0, len(products))
	for _, product := range products {
		dtos = append(dtos, ToProductDTO(product))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// Delete удаляет продукт
// @Summary Удалить продукт
// @Description Удаляет продукт по указанному идентификатору
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID продукта"
// @Success 204 "Успешно удалено"
// @Failure 400 {object} core.ErrorResponse "Неверный ID"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Продукт не найден"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/{id} [delete]
// @OperationId deleteProduct
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, productHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, productHandlerCode, "ИД параметр должен быть числвым! "+err.Error()))
		return
	}

	product, err := h.producterationService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}
	err = h.producterationService.Delete(ctx, product)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetByCategoryID возвращает продукты по категории
// @Summary Получить продукты по категории
// @Description Возвращает продукты по указанной категории
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id path int true "ID категории"
// @Success 200 {array} ProductDTO "Список продуктов"
// @Failure 400 {object} core.ErrorResponse "Неверная категория"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Продукты не найдены"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /products/category/{category_id} [get]
// @OperationId getProductByCategoryId
func (h *ProductHandler) GetByCategoryID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqCategoryID := r.PathValue("category_id")
	if reqCategoryID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, productHandlerCode, "Отсуствует ИД категории"))
		return
	}
	categoryID, err := strconv.Atoi(reqCategoryID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, productHandlerCode, "ИД категории должен быть числовым! "+err.Error()))
		return
	}

	products, err := h.producterationService.GetByCategoryID(ctx, uint(categoryID))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]ProductDTO, 0, len(products))
	for _, product := range products {
		dtos = append(dtos, ToProductDTO(product))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}
