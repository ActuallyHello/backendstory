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
