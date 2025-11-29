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
	cartItemHandlerCode = "CART_ITEM_HANDLER"
)

type CartItemHandler struct {
	validate        *validator.Validate
	cartItemService services.CartItemService
}

func NewCartItemHandler(
	cartItemService services.CartItemService,
) *CartItemHandler {
	return &CartItemHandler{
		validate:        validator.New(),
		cartItemService: cartItemService,
	}
}

// Create создает новый элемент корзины
// @Summary Создать элемент корзины
// @Description Создает новый элемент корзины
// @Tags CartItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CartItemCreateRequest true "Данные для создания элемента корзины"
// @Success 201 {object} dto.CartItemDTO "Созданный элемент корзины"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} dto.ErrorResponse "Элемент корзины уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/cart-items [post]
// @OperationId createCartItem
func (h *CartItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CartItemCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, cartItemHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, cartItemHandlerCode, err.Error()), details)
		return
	}

	cartItem := entities.CartItem{
		ProductID: req.ProductID,
		CartID:    req.CartID,
		Quantity:  req.Quantity,
	}
	cartItem, err := h.cartItemService.Create(ctx, cartItem)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToCartItemDTO(cartItem))
}

// Create создает новый элемент корзины
// @Summary Создать элемент корзины
// @Description Создает новый элемент корзины
// @Tags CartItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CartItemCreateRequest true "Данные для создания элемента корзины"
// @Success 201 {object} dto.CartItemDTO "Созданный элемент корзины"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} dto.ErrorResponse "Элемент корзины уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/cart-items [post]
// @OperationId createCartItem
func (h *CartItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CartItemUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, cartItemHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, cartItemHandlerCode, err.Error()), details)
		return
	}

	cartItem, err := h.cartItemService.GetByID(ctx, req.CartItemID)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	cartItem.Quantity = req.Quantity
	cartItem, err = h.cartItemService.Update(ctx, cartItem)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToCartItemDTO(cartItem))
}

// GetById возвращает элемент корзины по ID
// @Summary Получить элемент корзины по ID
// @Description Возвращает элемент корзины по указанному идентификатору
// @Tags CartItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID элемента корзины"
// @Success 200 {object} dto.CartItemDTO "Элемент корзины"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Элемент корзины не найден"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/cart-items/{id} [get]
// @OperationId getCartItemById
func (h *CartItemHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, cartItemHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, cartItemHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	cartItem, err := h.cartItemService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ToCartItemDTO(cartItem))
}

// GetWithSearchCriteria выполняет поиск элементов корзины по критериям
// @Summary Поиск элементов корзины
// @Description Выполняет поиск элементов корзины по заданным критериям
// @Tags CartItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SearchCriteria true "Критерии поиска"
// @Success 200 {array} dto.CartItemDTO "Список найденных элементов корзины"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/cart-items/search [post]
// @OperationId searchCartItem
func (h *CartItemHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, cartItemHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, cartItemHandlerCode, err.Error()), details)
		return
	}

	cartItems, err := h.cartItemService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.CartItemDTO, 0, len(cartItems))
	for _, cartItem := range cartItems {
		dtos = append(dtos, dto.ToCartItemDTO(cartItem))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetByCartID возвращает элементы корзины по ID корзины
// @Summary Получить элементы корзины по ID корзины
// @Description Возвращает все элементы корзины по указанному идентификатору корзины
// @Tags CartItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cart_id path int true "ID корзины"
// @Success 200 {array} dto.CartItemDTO "Список элементов корзины"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID корзины"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Элементы корзины не найдены"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/cart-items/cart/{cart_id} [get]
// @OperationId getCartItemByCartId
func (h *CartItemHandler) GetByCartID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("cart_id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, cartItemHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, cartItemHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	cartItems, err := h.cartItemService.GetByCartID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.CartItemDTO, 0, len(cartItems))
	for _, cartItem := range cartItems {
		dtos = append(dtos, dto.ToCartItemDTO(cartItem))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}
