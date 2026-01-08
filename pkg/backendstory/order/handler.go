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

// Create создает новый заказ
// @Summary Создать заказ
// @Description Создает новый заказ на основе товаров из корзины
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body OrderCreateRequest true "Данные для создания заказа"
// @Success 201 {object} OrderDTO "Созданный заказ"
// @Failure 400 {object} core.ErrorResponse "Неверный запрос"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} core.ErrorResponse "Конфликт (заказ уже существует)"
// @Failure 422 {object} core.ValidationErrorResponse "Ошибка валидации"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders [post]
// @Id createOrder
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

// ChangeStatus изменяет статус заказа
// @Summary Изменить статус заказа
// @Description Изменяет статус заказа и назначает менеджера (если не назначен)
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заказа"
// @Param status path string true "Новый статус заказа (pending, processing, shipped, delivered, cancelled)"
// @Success 200 {object} OrderDTO "Обновленный заказ"
// @Failure 400 {object} core.ErrorResponse "Неверный запрос или статус"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Заказ не найден"
// @Failure 409 {object} core.ErrorResponse "Конфликт статусов"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders/{id}/status/{status} [patch]
// @Id changeOrderStatus
func (h *OrderHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсутствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}

	status := r.PathValue("status")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсутствует действие к заказу"))
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

// Delete удаляет заказ
// @Summary Удалить заказ
// @Description Удаляет заказ по ID (только для заказов в определенных статусах)
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заказа"
// @Success 204 "Заказ успешно удален"
// @Failure 400 {object} core.ErrorResponse "Неверный ID заказа"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Заказ не найден"
// @Failure 409 {object} core.ErrorResponse "Нельзя удалить заказ в текущем статусе"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders/{id} [delete]
// @Id deleteOrder
func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсутствует ИД параметр"))
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

// GetById возвращает заказ по ID
// @Summary Получить заказ по ID
// @Description Возвращает заказ по указанному идентификатору
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заказа"
// @Success 200 {object} OrderDTO "Заказ"
// @Failure 400 {object} core.ErrorResponse "Неверный ID заказа"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Заказ не найден"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders/{id} [get]
// @Id getOrderById
func (h *OrderHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderHandlerCode, "Отсутствует ИД параметр"))
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

// GetWithSearchCriteria возвращает список заказов по критериям поиска
// @Summary Поиск заказов
// @Description Возвращает список заказов по указанным критериям поиска с пагинацией
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.SearchCriteria true "Критерии поиска"
// @Success 200 {array} OrderDTO "Список заказов"
// @Failure 400 {object} core.ErrorResponse "Неверные критерии поиска"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 422 {object} core.ValidationErrorResponse "Ошибка валидации"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders/search [post]
// @Id searchOrders
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

// GetByStatus возвращает заказы по статусу
// @Summary Получить заказы по статусу
// @Description Возвращает список заказов с указанным статусом
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status path string true "Статус заказа (pending, processing, shipped, delivered, cancelled)"
// @Success 200 {array} OrderDTO "Список заказов"
// @Failure 400 {object} core.ErrorResponse "Неверный статус"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders/status/{status} [get]
// @Id getOrdersByStatus
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

// GetByClientID возвращает заказы по ID клиента
// @Summary Получить заказы клиента
// @Description Возвращает список заказов по ID клиента
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param client_id path int true "ID клиента"
// @Success 200 {array} OrderDTO "Список заказов клиента"
// @Failure 400 {object} core.ErrorResponse "Неверный ID клиента"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders/client/{client_id} [get]
// @Id getOrdersByClientId
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

// GetByManagerID возвращает заказы по ID менеджера
// @Summary Получить заказы менеджера
// @Description Возвращает список заказов, назначенных на менеджера
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param manager_id path int true "ID менеджера"
// @Success 200 {array} OrderDTO "Список заказов менеджера"
// @Failure 400 {object} core.ErrorResponse "Неверный ID менеджера"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders/manager/{manager_id} [get]
// @Id getOrdersByManagerId
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

// GetByManagerIDAndStatus возвращает заказы по ID менеджера и статусу
// @Summary Получить заказы менеджера по статусу
// @Description Возвращает список заказов, назначенных на менеджера с указанным статусом
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param manager_id path int true "ID менеджера"
// @Param status path string true "Статус заказа (pending, processing, shipped, delivered, cancelled)"
// @Success 200 {array} OrderDTO "Список заказов"
// @Failure 400 {object} core.ErrorResponse "Неверные параметры"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /orders/manager/{manager_id}/status/{status} [get]
// @Id getOrdersByManagerIdAndStatus
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
