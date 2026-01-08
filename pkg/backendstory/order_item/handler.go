package orderitem

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
)

const (
	orderItemHandlerCode = "ORDER_ITEM_HANDLER"
)

type OrderItemHandler struct {
	validate         *validator.Validate
	orderItemService OrderItemService
}

func NewOrderItemHandler(
	orderItemService OrderItemService,
) *OrderItemHandler {
	return &OrderItemHandler{
		validate:         validator.New(),
		orderItemService: orderItemService,
	}
}

func (h *OrderItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req OrderItemCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, orderItemHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, orderItemHandlerCode, err.Error()), details)
		return
	}

	orderItem := OrderItem{
		OrderID:    req.OrderID,
		CartItemID: req.CartItemID,
	}
	orderItem, err := h.orderItemService.Create(ctx, orderItem)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToOrderItemDTO(orderItem))
}

func (h *OrderItemHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderItemHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderItemHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}
	status := r.PathValue("status")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderItemHandlerCode, "Отсуствует действие к заказу"))
		return
	}
	if status == "" {
		core.HandleError(w, r, core.NewLogicalError(err, orderItemHandlerCode, "Параметр действия над заказом пустой!"))
		return
	}

	orderItem, err := h.orderItemService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(nil, orderItemHandlerCode, "Элемента заказа не существует"))
		return
	}

	orderItem, err = h.orderItemService.ChangeStatus(ctx, orderItem, status)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToOrderItemDTO(orderItem))
}

func (h *OrderItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderItemHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderItemHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	orderItem, err := h.orderItemService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(nil, orderItemHandlerCode, "Элемента заказа не существует"))
		return
	}

	err = h.orderItemService.Delete(ctx, orderItem)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderItemHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderItemHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderItemHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	orderItem, err := h.orderItemService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToOrderItemDTO(orderItem))
}

func (h *OrderItemHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req core.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, orderItemHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, orderItemHandlerCode, err.Error()), details)
		return
	}

	orderItems, err := h.orderItemService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]OrderItemDTO, 0, len(orderItems))
	for _, orderItem := range orderItems {
		dtos = append(dtos, ToOrderItemDTO(orderItem))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

func (h *OrderItemHandler) GetByOrderID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqOrderID := r.PathValue("order_id")
	if reqOrderID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderItemHandlerCode, "Отсутствует ИД параметр"))
		return
	}
	orderID, err := strconv.Atoi(reqOrderID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderItemHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	orderItems, err := h.orderItemService.GetByOrderID(ctx, uint(orderID))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	orderItemDTOs := make([]OrderItemDTO, 0, len(orderItems))
	for _, orderItem := range orderItems {
		orderItemDTOs = append(orderItemDTOs, ToOrderItemDTO(orderItem))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderItemDTOs)
}
