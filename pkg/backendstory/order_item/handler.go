package orderitem

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
)

const (
	orderItemHandlerCode = "ORDER_ITEM_HANDLER"
)

type OrderItemHandler struct {
	validate         *validator.Validate
	orderItemService OrderItemService
	enumValueService enumvalue.EnumValueService
}

func NewOrderItemHandler(
	orderItemService OrderItemService,
	enumValueService enumvalue.EnumValueService,
) *OrderItemHandler {
	return &OrderItemHandler{
		validate:         validator.New(),
		orderItemService: orderItemService,
		enumValueService: enumValueService,
	}
}

// Create создает новый элемент заказа
// @Summary Создать элемент заказа
// @Description Создает новый элемент заказа, связывая товар из корзины с заказом
// @Tags OrderItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body OrderItemCreateRequest true "Данные для создания элемента заказа"
// @Success 201 {object} OrderItemDTO "Созданный элемент заказа"
// @Failure 400 {object} core.ErrorResponse "Неверный запрос"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 409 {object} core.ErrorResponse "Конфликт (элемент заказа уже существует)"
// @Failure 422 {object} core.ValidationErrorResponse "Ошибка валидации"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /order-items [post]
// @Id createOrderItem
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

	orderItemStatus, err := h.enumValueService.GetByID(ctx, orderItem.StatusID)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToOrderItemDTO(orderItem, enumvalue.ToEnumValueDTO(orderItemStatus)))
}

// ChangeStatus изменяет статус элемента заказа
// @Summary Изменить статус элемента заказа
// @Description Изменяет статус элемента заказа (например, готов к отгрузке, отгружен и т.д.)
// @Tags OrderItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID элемента заказа"
// @Param status path string true "Новый статус элемента заказа (pending, ready_to_ship, shipped, delivered, cancelled)"
// @Success 201 {object} OrderItemDTO "Обновленный элемент заказа"
// @Failure 400 {object} core.ErrorResponse "Неверный запрос или статус"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Элемент заказа не найден"
// @Failure 409 {object} core.ErrorResponse "Конфликт статусов"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /order-items/{id}/change-status/{status} [post]
// @Id changeOrderItemStatus
func (h *OrderItemHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderItemHandlerCode, "Отсутствует ИД параметр"))
		return
	}
	id, err := strconv.Atoi(reqID)
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, orderItemHandlerCode, "ИД параметр должен быть числовым!"+err.Error()))
		return
	}
	status := r.PathValue("status")
	if status == "" {
		core.HandleError(w, r, core.NewLogicalError(err, orderItemHandlerCode, "Параметр действия над элементом заказа пустой!"))
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

	orderItemStatus, err := h.enumValueService.GetByID(ctx, orderItem.StatusID)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToOrderItemDTO(orderItem, enumvalue.ToEnumValueDTO(orderItemStatus)))
}

// Delete удаляет элемент заказа
// @Summary Удалить элемент заказа
// @Description Удаляет элемент заказа по ID (только для элементов в определенных статусах)
// @Tags OrderItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID элемента заказа"
// @Success 204 "Элемент заказа успешно удален"
// @Failure 400 {object} core.ErrorResponse "Неверный ID элемента заказа"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Элемент заказа не найден"
// @Failure 409 {object} core.ErrorResponse "Нельзя удалить элемент заказа в текущем статусе"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /order-items/{id} [delete]
// @Id deleteOrderItem
func (h *OrderItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderItemHandlerCode, "Отсутствует ИД параметр"))
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

// GetById возвращает элемент заказа по ID
// @Summary Получить элемент заказа по ID
// @Description Возвращает элемент заказа по указанному идентификатору
// @Tags OrderItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID элемента заказа"
// @Success 200 {object} OrderItemDTO "Элемент заказа"
// @Failure 400 {object} core.ErrorResponse "Неверный ID элемента заказа"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 404 {object} core.ErrorResponse "Элемент заказа не найден"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /order-items/{id} [get]
// @Id getOrderItemById
func (h *OrderItemHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqID := r.PathValue("id")
	if reqID == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, orderItemHandlerCode, "Отсутствует ИД параметр"))
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

	orderItemStatus, err := h.enumValueService.GetByID(ctx, orderItem.StatusID)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToOrderItemDTO(orderItem, enumvalue.ToEnumValueDTO(orderItemStatus)))
}

// GetWithSearchCriteria возвращает список элементов заказа по критериям поиска
// @Summary Поиск элементов заказа
// @Description Возвращает список элементов заказа по указанным критериям поиска с пагинацией
// @Tags OrderItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body core.SearchCriteria true "Критерии поиска"
// @Success 200 {array} OrderItemDTO "Список элементов заказа"
// @Failure 400 {object} core.ErrorResponse "Неверные критерии поиска"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 422 {object} core.ValidationErrorResponse "Ошибка валидации"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /order-items/search [post]
// @Id searchOrderItems
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
		orderItemStatus, err := h.enumValueService.GetByID(ctx, orderItem.StatusID)
		if err != nil {
			core.HandleError(w, r, err)
			return
		}
		dtos = append(dtos, ToOrderItemDTO(orderItem, enumvalue.ToEnumValueDTO(orderItemStatus)))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// GetByOrderID возвращает элементы заказа по ID заказа
// @Summary Получить элементы заказа
// @Description Возвращает список элементов заказа по ID заказа
// @Tags OrderItems
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id path int true "ID заказа"
// @Success 200 {array} OrderItemDTO "Список элементов заказа"
// @Failure 400 {object} core.ErrorResponse "Неверный ID заказа"
// @Failure 401 {object} core.ErrorResponse "Не авторизован"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /order-items/order/{order_id} [get]
// @Id getOrderItemsByOrderId
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
		orderItemStatus, err := h.enumValueService.GetByID(ctx, orderItem.StatusID)
		if err != nil {
			core.HandleError(w, r, err)
			return
		}
		orderItemDTOs = append(orderItemDTOs, ToOrderItemDTO(orderItem, enumvalue.ToEnumValueDTO(orderItemStatus)))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderItemDTOs)
}
