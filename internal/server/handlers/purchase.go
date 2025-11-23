package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	appErr "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/ActuallyHello/backendstory/internal/server/handlers/common"
	"github.com/ActuallyHello/backendstory/internal/server/middleware"
	"github.com/ActuallyHello/backendstory/internal/services"
	"github.com/go-playground/validator/v10"
)

const (
	purchaseHandlerCode = "PURCHASE_HANDLER"
)

type PurchaseHandler struct {
	validate *validator.Validate

	purchaseService services.PurchaseService
	productService  services.ProductService
	personService   services.PersonService
}

func NewPurchaseHandler(
	purchaseService services.PurchaseService,
) *PurchaseHandler {
	return &PurchaseHandler{
		validate:        validator.New(),
		purchaseService: purchaseService,
	}
}

// AddToCart добавляет товар в корзину
// @Summary Добавить товар в корзину
// @Description Добавляет указанное количество товара в корзину пользователя
// @Tags Purchases
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddToCartRequest true "Данные для добавления в корзину"
// @Success 201 "Товар успешно добавлен в корзину"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Товар или пользователь не найден"
// @Failure 409 {object} dto.ErrorResponse "Недостаточно товара на складе"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/purchases/cart [post]
// @OperationId addToCart
func (h *PurchaseHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, purchaseHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, purchaseHandlerCode, err.Error()), details)
		return
	}

	userInfo, err := middleware.GetUserInfoCtx(ctx)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}
	person, err := h.personService.GetByUserLogin(ctx, userInfo.Username)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, purchaseHandlerCode, "Person doesn't exists with such username"+userInfo.Username))
		return
	}
	product, err := h.productService.GetByID(ctx, req.ProductID)
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, purchaseHandlerCode, fmt.Sprintf("Product with id=%d doesn't exists", req.ProductID)))
		return
	}

	if err := h.purchaseService.AddToCart(ctx, product, person, req.Quantity); err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, purchaseHandlerCode, fmt.Sprintf("Couldn't purchase product %s with %d items", product.Code, req.Quantity)))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Purchase оформляет покупку корзины
// @Summary Оформить покупку корзины
// @Description Оформляет покупку всех товаров в указанной корзине
// @Tags Purchases
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.PurchaseRequest true "Данные для оформления покупки"
// @Success 201 "Покупка успешно оформлена"
// @Failure 400 {object} dto.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} dto.ErrorResponse "Не авторизован"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} dto.ErrorResponse "Корзина не найдена"
// @Failure 409 {object} dto.ErrorResponse "Недостаточно товаров на складе"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/purchases [post]
// @OperationId purchaseCart
func (h *PurchaseHandler) Purchase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.HandleError(w, r, appErr.NewTechnicalError(err, purchaseHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := common.CollectValidationDetails(err)
		middleware.HandleValidationError(w, r, appErr.NewLogicalError(err, purchaseHandlerCode, err.Error()), details)
		return
	}

	if err := h.purchaseService.Purchase(ctx, req.CartID); err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, purchaseHandlerCode, "Невозможно оформить корзину!"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
