package cartitem

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
)

const (
	cartItemHandlerCode = "CART_ITEM_HANDLER"
)

type CartItemHandler struct {
	validate        *validator.Validate
	cartItemService CartItemService
}

func NewCartItemHandler(
	cartItemService CartItemService,
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
// @Param request body CartItemCreateRequest true "Данные для создания элемента корзины"
// @Success 201 {object} CartItemDTO "Созданный элемент корзины"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} core.ErrorResponse "Элемент корзины уже существует"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /cart-items [post]
// @ID createCartItem
func (h *CartItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CartItemCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, cartItemHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, cartItemHandlerCode, err.Error()), details)
		return
	}

	cartItem := CartItem{
		ProductID: req.ProductID,
		CartID:    req.CartID,
		Quantity:  req.Quantity,
	}
	cartItem, err := h.cartItemService.Create(ctx, cartItem)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToCartItemDTO(cartItem))
}

// Create создает новый элемент корзины
// @Summary Создать элемент корзины
// @Description Создает новый элемент корзины
// @Tags CartItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CartItemCreateRequest true "Данные для создания элемента корзины"
// @Success 201 {object} CartItemDTO "Созданный элемент корзины"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} core.ErrorResponse "Элемент корзины уже существует"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /cart-items [put]
// @ID updateCartItem
func (h *CartItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CartItemUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, cartItemHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, cartItemHandlerCode, err.Error()), details)
		return
	}

	cartItem, err := h.cartItemService.GetByID(ctx, req.CartItemID)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	cartItem.Quantity = req.Quantity
	cartItem, err = h.cartItemService.Update(ctx, cartItem)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToCartItemDTO(cartItem))
}

// GetById возвращает элемент корзины по ID
// @Summary Получить элемент корзины по ID
// @Description Возвращает элемент корзины по указанному идентификатору
// @Tags CartItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID элемента корзины"
// @Success 200 {object} CartItemDTO "Элемент корзины"
// @Failure 400 {object} core.ErrorResponse "Неверный ID"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Элемент корзины не найден"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /cart-items/{id} [get]
// @ID getCartItemById
func (h *CartItemHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, cartItemHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, cartItemHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	cartItem, err := h.cartItemService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToCartItemDTO(cartItem))
}

// GetWithSearchCriteria выполняет поиск элементов корзины по критериям
// @Summary Поиск элементов корзины
// @Description Выполняет поиск элементов корзины по заданным критериям
// @Tags CartItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.SearchCriteria true "Критерии поиска"
// @Success 200 {array} CartItemDTO "Список найденных элементов корзины"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /cart-items/search [post]
// @ID getCartItemSearch
func (h *CartItemHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req core.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, cartItemHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, cartItemHandlerCode, err.Error()), details)
		return
	}

	cartItems, err := h.cartItemService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]CartItemDTO, 0, len(cartItems))
	for _, cartItem := range cartItems {
		dtos = append(dtos, ToCartItemDTO(cartItem))
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
// @Success 200 {array} CartItemDTO "Список элементов корзины"
// @Failure 400 {object} core.ErrorResponse "Неверный ID корзины"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Элементы корзины не найдены"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /cart-items/cart/{cart_id} [get]
// @ID getCartItemByCartId
func (h *CartItemHandler) GetByCartID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("cart_id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, cartItemHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, cartItemHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	cartItems, err := h.cartItemService.GetByCartID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]CartItemDTO, 0, len(cartItems))
	for _, cartItem := range cartItems {
		dtos = append(dtos, ToCartItemDTO(cartItem))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}
