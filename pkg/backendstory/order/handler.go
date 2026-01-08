package order

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/auth"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/person"
	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
)

const (
	orderHandlerCode = "ORDER_HANDLER"
)

type OrderHandler struct {
	validate      *validator.Validate
	orderService  OrderService
	personService person.PersonService
}

func NewOrderHandler(
	orderService OrderService,
	personService person.PersonService,
) *OrderHandler {
	return &OrderHandler{
		validate:      validator.New(),
		orderService:  orderService,
		personService: personService,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req OrderCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, orderHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, orderHandlerCode, err.Error()), details)
		return
	}

	order := Order{
		ClientID: req.ClientID,
	}
	order, err := h.orderService.Create(ctx, order, req.CartItemIDs)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToOrderDTO(order))
}

func (h *OrderHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	status := r.PathValue("status")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсуствует действие к заказу"))
		return
	}
	if status == "" {
		core.HandleError(w, r, core.NewLogicalError(err, orderHandlerCode, "Параметр действия над заказом пустой!"))
		return
	}

	userinfo, err := auth.GetUserInfoCtx(ctx)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}
	manager, err := h.personService.GetByUserLogin(ctx, userinfo.Username)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	order, err := h.orderService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	order.ManagerID = sql.NullInt32{
		Int32: int32(manager.ID),
		Valid: true,
	}
	order, err = h.orderService.ChangeStatus(ctx, order, status)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToOrderDTO(order))
}

func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	order, err := h.orderService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	err = h.orderService.Delete(ctx, order)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсуствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	order, err := h.orderService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToOrderDTO(order))
}

func (h *OrderHandler) GetWithSearchCriteria(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req core.SearchCriteria
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		core.HandleError(w, r, core.NewTechnicalError(err, orderHandlerCode, err.Error()))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		details := core.CollectValidationDetails(err)
		core.HandleValidationError(w, r, core.NewLogicalError(err, orderHandlerCode, err.Error()), details)
		return
	}

	orders, err := h.orderService.GetWithSearchCriteria(ctx, req)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]OrderDTO, 0, len(orders))
	for _, order := range orders {
		dtos = append(dtos, ToOrderDTO(order))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

func (h *OrderHandler) GetByStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status := r.PathValue("status")
	if status == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсутствует параметр статус"))
		return
	}

	orders, err := h.orderService.GetByStatus(ctx, status)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	orderDTOs := make([]OrderDTO, 0, len(orders))
	for _, order := range orders {
		orderDTOs = append(orderDTOs, ToOrderDTO(order))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderDTOs)
}

func (h *OrderHandler) GetByClientID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqClientID := r.PathValue("client_id")
	if reqClientID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсутствует ИД параметр"))
		return
	}
	clientID, err := strconv.Atoi(reqClientID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	orders, err := h.orderService.GetByClientID(ctx, uint(clientID))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	orderDTOs := make([]OrderDTO, 0, len(orders))
	for _, order := range orders {
		orderDTOs = append(orderDTOs, ToOrderDTO(order))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderDTOs)
}

func (h *OrderHandler) GetByManagerID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqManagerID := r.PathValue("manager_id")
	if reqManagerID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсутствует ИД параметр"))
		return
	}
	managerID, err := strconv.Atoi(reqManagerID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	orders, err := h.orderService.GetByManagerID(ctx, uint(managerID))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	orderDTOs := make([]OrderDTO, 0, len(orders))
	for _, order := range orders {
		orderDTOs = append(orderDTOs, ToOrderDTO(order))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderDTOs)
}

func (h *OrderHandler) GetByManagerIDAndStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqManagerID := r.PathValue("manager_id")
	if reqManagerID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсутствует ИД параметр"))
		return
	}
	managerID, err := strconv.Atoi(reqManagerID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}
	status := r.PathValue("status")
	if status == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсутствует параметр статус"))
		return
	}

	orders, err := h.orderService.GetByManagerIDAndStatus(ctx, uint(managerID), status)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	orderDTOs := make([]OrderDTO, 0, len(orders))
	for _, order := range orders {
		orderDTOs = append(orderDTOs, ToOrderDTO(order))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderDTOs)
}
