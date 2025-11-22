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
	cartHandlerCode = "CART_HANDLER"
)

type CartHandler struct {
	validate    *validator.Validate
	cartService services.CartService
}

func NewCartHandler(
	cartService services.CartService,
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
// @Param request body dto.CartCreateRequest true "Данные для создания корзины"
// @Success 201 {object} dto.CartDTO "Созданная корзина"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} dto.ErrorResponse "Корзина для пользователя уже существует"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/carts [post]
func (h *CartHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CartCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, cartHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, cartHandlerCode, err.Error()), details)
		return
	}

	cart := entities.Cart{
		PersonID: req.PersonID,
	}
	cart, err := h.cartService.Create(ctx, cart)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToCartDTO(cart))
}

// GetById возвращает корзину по ID
// @Summary Получить корзину по ID
// @Description Возвращает корзину по указанному идентификатору
// @Tags Carts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID корзины"
// @Success 200 {object} dto.CartDTO "Корзина"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Корзина не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/carts/{id} [get]
func (h *CartHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, cartHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, cartHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	cart, err := h.cartService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ToCartDTO(cart))
}

// GetWithSearchCriteria выполняет поиск корзин по критериям
// @Summary Поиск корзин
// @Description Выполняет поиск корзин по заданным критериям
// @Tags Carts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SearchCriteria true "Критерии поиска"
// @Success 200 {array} dto.CartDTO "Список найденных корзин"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/carts/search [post]
func (h *CartHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, cartHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, cartHandlerCode, err.Error()), details)
		return
	}

	carts, err := h.cartService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.CartDTO, 0, len(carts))
	for _, cart := range carts {
		dtos = append(dtos, dto.ToCartDTO(cart))
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
// @Success 200 {object} dto.CartDTO "Корзина"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID пользователя"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Корзина не найдена"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/carts/person/{person_id} [get]
func (h *CartHandler) GetByPersonID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("person_id")
	if reqID == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, cartHandlerCode, "ID parameter missing"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, cartHandlerCode, "ID parameter must be integer!"+err.Error()))
		return
	}

	cart, err := h.cartService.GetByPersonID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ToCartDTO(cart))
}
