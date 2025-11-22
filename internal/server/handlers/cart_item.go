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
	}
	cartItem, err := h.cartItemService.Create(ctx, cartItem)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToCartItemDTO(cartItem))
}

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

func (h *CartItemHandler) GetByPersonID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("person_id")
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

	dtos := make([]dto.CartItemDTO, len(cartItems))
	for _, cartItem := range cartItems {
		dtos = append(dtos, dto.ToCartItemDTO(cartItem))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}
