package purchase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/auth"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/person"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/product"
	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
)

const (
	purchaseHandlerCode = "PURCHASE_HANDLER"
)

type PurchaseHandler struct {
	validate *validator.Validate

	purchaseService PurchaseService
	productService  product.ProductService
	personService   person.PersonService
}

func NewPurchaseHandler(
	purchaseService PurchaseService,
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
// @Param request body AddToCartRequest true "Данные для добавления в корзину"
// @Success 201 "Товар успешно добавлен в корзину"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 403 {object} ErrorResponse "Доступ запрещен"
// @Failure 404 {object} ErrorResponse "Товар или пользователь не найден"
// @Failure 409 {object} ErrorResponse "Недостаточно товара на складе"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/purchases/cart [post]
// @OperationId addToCart
func (h *PurchaseHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, purchaseHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, purchaseHandlerCode, err.Error()), details)
		return
	}

	userInfo, err := auth.GetUserInfoCtx(ctx)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}
	person, err := h.personService.GetByUserLogin(ctx, userInfo.Username)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, purchaseHandlerCode, "Person doesn't exists with such username"+userInfo.Username))
		return
	}
	product, err := h.productService.GetByID(ctx, req.ProductID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, purchaseHandlerCode, fmt.Sprintf("Product with id=%d doesn't exists", req.ProductID)))
		return
	}

	if err := h.purchaseService.AddToCart(ctx, product, person, req.Quantity); err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, purchaseHandlerCode, fmt.Sprintf("Couldn't purchase product %s with %d items", product.Code, req.Quantity)))
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
// @Param request body PurchaseRequest true "Данные для оформления покупки"
// @Success 201 "Покупка успешно оформлена"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 401 {object} ErrorResponse "Не авторизован"
// @Failure 403 {object} ErrorResponse "Доступ запрещен"
// @Failure 404 {object} ErrorResponse "Корзина не найдена"
// @Failure 409 {object} ErrorResponse "Недостаточно товаров на складе"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/purchases [post]
// @OperationId purchaseCart
func (h *PurchaseHandler) Purchase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, purchaseHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, purchaseHandlerCode, err.Error()), details)
		return
	}

	if err := h.purchaseService.Purchase(ctx, req.CartID); err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, purchaseHandlerCode, "Невозможно оформить корзину!"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
