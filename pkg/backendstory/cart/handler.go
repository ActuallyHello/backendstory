package cart

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
)

const (
	cartHandlerCode = "CART_HANDLER"
)

type CartHandler struct {
	validate    *validator.Validate
	cartService CartService
}

func NewCartHandler(
	cartService CartService,
) *CartHandler {
	return &CartHandler{
		validate:    validator.New(),
		cartService: cartService,
	}
}

// Create создает новую корзину
// @Summary Создать корзину
// @Description Создает новую корзину для пользователя
// @Tags Carts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CartCreateRequest true "Данные для создания корзины"
// @Success 201 {object} CartDTO "Созданная корзина"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} core.ErrorResponse "Корзина для пользователя уже существует"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /carts [post]
// @OperationId createCart
func (h *CartHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CartCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, cartHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, cartHandlerCode, err.Error()), details)
		return
	}

	cart := Cart{
		PersonID: req.PersonID,
	}
	cart, err := h.cartService.Create(ctx, cart)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToCartDTO(cart))
}

// GetById возвращает корзину по ID
// @Summary Получить корзину по ID
// @Description Возвращает корзину по указанному идентификатору
// @Tags Carts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID корзины"
// @Success 200 {object} CartDTO "Корзина"
// @Failure 400 {object} core.ErrorResponse "Неверный ID"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Корзина не найдена"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /carts/{id} [get]
// @OperationId GetCartById
func (h *CartHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, cartHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, cartHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	cart, err := h.cartService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToCartDTO(cart))
}

// GetWithSearchCriteria выполняет поиск корзин по критериям
// @Summary Поиск корзин
// @Description Выполняет поиск корзин по заданным критериям
// @Tags Carts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.SearchCriteria true "Критерии поиска"
// @Success 200 {array} CartDTO "Список найденных корзин"
// @Failure 400 {object} core.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /carts/search [post]
// @OperationId searchCart
func (h *CartHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req core.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, cartHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, cartHandlerCode, err.Error()), details)
		return
	}

	carts, err := h.cartService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]CartDTO, 0, len(carts))
	for _, cart := range carts {
		dtos = append(dtos, ToCartDTO(cart))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetByPersonID возвращает корзину по ID пользователя
// @Summary Получить корзину по ID пользователя
// @Description Возвращает корзину по указанному идентификатору пользователя
// @Tags Carts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param person_id path int true "ID пользователя"
// @Success 200 {object} CartDTO "Корзина"
// @Failure 400 {object} core.ErrorResponse "Неверный ID пользователя"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Корзина не найдена"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /carts/person/{person_id} [get]
// @OperationId getCartByPersonId
func (h *CartHandler) GetByPersonID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("person_id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, cartHandlerCode, "Отсутствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, cartHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	cart, err := h.cartService.GetByPersonID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToCartDTO(cart))
}
