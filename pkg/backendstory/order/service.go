package order

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/enum"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
	orderitem "github.com/ActuallyHello/backendstory/pkg/backendstory/order_item"
	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	orderServiceCode = "ORDER_SERVICE"
)

type OrderService interface {
	core.BaseService[Order]

	Create(ctx context.Context, order Order, cartItemIDs []uint) (Order, error)
	Update(ctx context.Context, order Order) (Order, error)
	Delete(ctx context.Context, order Order) error

	Approve(ctx context.Context, order Order) (Order, error)
	Cancel(ctx context.Context, order Order) (Order, error)
	ChangeStatus(ctx context.Context, order Order, status string) (Order, error)

	GetByStatus(ctx context.Context, status string) ([]Order, error)
	GetByClientID(ctx context.Context, clientID uint) ([]Order, error)
	GetByManagerID(ctx context.Context, managerID uint) ([]Order, error)
	GetByManagerIDAndStatus(ctx context.Context, managerID uint, status string) ([]Order, error)
}

type orderService struct {
	core.BaseServiceImpl[Order]
	orderRepo        OrderRepository
	txManager        core.TxManager
	enumService      enum.EnumService
	enumValueService enumvalue.EnumValueService
	orderItemService orderitem.OrderItemService
}

func NewOrderService(
	orderRepo OrderRepository,
	txManager core.TxManager,
	enumService enum.EnumService,
	enumValueService enumvalue.EnumValueService,
	orderItemService orderitem.OrderItemService,
) *orderService {
	return &orderService{
		BaseServiceImpl:  *core.NewBaseServiceImpl(orderRepo),
		orderRepo:        orderRepo,
		txManager:        txManager,
		enumService:      enumService,
		enumValueService: enumValueService,
		orderItemService: orderItemService,
	}
}

func (s *orderService) Create(ctx context.Context, order Order, cartItemIDs []uint) (Order, error) {
	var newOrder Order
	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		pendingOrderStatus, err := s.enumValueService.GetByCodeAndEnumCode(ctx, PendingOrderStatus, OrderStatus)
		if err != nil {
			return err
		}

		order.StatusID = pendingOrderStatus.ID
		order, err := s.orderRepo.Create(ctx, order)
		if err != nil {
			return err
		}
		newOrder = order

		for _, cartItemID := range cartItemIDs {
			if _, err := s.orderItemService.Create(ctx, orderitem.OrderItem{
				OrderID:    order.ID,
				CartItemID: cartItemID,
			}); err != nil {
				return err
			}
		}

		return nil
	})
	return newOrder, err
}

func (s *orderService) ChangeStatus(ctx context.Context, order Order, status string) (Order, error) {
	switch status {
	case ApprovedOrderStatus:
		order, err := s.Approve(ctx, order)
		if err != nil {
			return Order{}, err
		}
		return order, nil
	case CancelledOrderStatus:
		order, err := s.Cancel(ctx, order)
		if err != nil {
			return Order{}, err
		}
		return order, nil
	default:
		return Order{}, core.NewLogicalError(nil, orderServiceCode, "Неизвестный статус заказа!")
	}
}

func (s *orderService) Approve(ctx context.Context, order Order) (Order, error) {
	var approvedOrder Order
	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		currentOrderStatus, err := s.enumValueService.GetByID(ctx, order.StatusID)
		if currentOrderStatus.Code == CancelledOrderStatus {
			return core.NewLogicalError(nil, orderServiceCode, "Невозможно подтвердить отмененный заказ!")
		}

		orderItems, err := s.orderItemService.GetByOrderID(ctx, order.ID)
		if err != nil {
			return err
		}

		for _, orderItem := range orderItems {
			if _, err := s.orderItemService.Approve(ctx, orderItem); err != nil {
				return err
			}
		}

		approvedStatus, err := s.enumValueService.GetByCodeAndEnumCode(ctx, ApprovedOrderStatus, OrderStatus)
		if err != nil {
			return err
		}

		order.StatusID = approvedStatus.ID
		order, err = s.orderRepo.Update(ctx, order)
		if err != nil {
			return err
		}
		approvedOrder = order

		return nil
	})

	return approvedOrder, err
}

func (s *orderService) Cancel(ctx context.Context, order Order) (Order, error) {
	var cancelledOrder Order
	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		orderItems, err := s.orderItemService.GetByOrderID(ctx, order.ID)
		if err != nil {
			return err
		}
		for _, orderItem := range orderItems {
			if _, err := s.orderItemService.Cancel(ctx, orderItem); err != nil {
				return err
			}
		}

		cancelled, err := s.enumValueService.GetByCodeAndEnumCode(ctx, CancelledOrderStatus, OrderStatus)
		if err != nil {
			return err
		}

		order.StatusID = cancelled.ID
		order, err = s.orderRepo.Update(ctx, order)
		if err != nil {
			return err
		}
		cancelledOrder = order

		return nil
	})

	return cancelledOrder, err
}

func (s *orderService) Update(ctx context.Context, order Order) (Order, error) {
	order, err := s.GetRepo().Update(ctx, order)
	if err != nil {
		return Order{}, core.NewTechnicalError(err, orderServiceCode, "Ошибка при обновлении заказа!")
	}
	return order, nil
}

func (s *orderService) Delete(ctx context.Context, order Order) error {
	return s.txManager.Do(ctx, func(ctx context.Context) error {
		orderItems, err := s.orderItemService.GetByOrderID(ctx, order.ID)
		if err != nil {
			return err
		}
		for _, orderItem := range orderItems {
			if err := s.orderItemService.Delete(ctx, orderItem); err != nil {
				return err
			}
		}

		err = s.GetRepo().Delete(ctx, order)
		if err != nil {
			return core.NewTechnicalError(err, orderServiceCode, "Ошибка при удалении заказа")
		}

		return nil
	})
}

func (s *orderService) GetByClientID(ctx context.Context, clientID uint) ([]Order, error) {
	orders, err := s.orderRepo.FindByClientID(ctx, clientID)
	if err != nil {
		if errors.Is(err, &core.NotFoundError{}) {
			return nil, core.NewLogicalError(err, orderServiceCode, err.Error())
		}
		return nil, core.NewTechnicalError(err, orderServiceCode, "Ошибка при поиске заказа у клиента")
	}
	return orders, nil
}

func (s *orderService) GetByStatus(ctx context.Context, status string) ([]Order, error) {
	orderStatus, err := s.enumValueService.GetByCodeAndEnumCode(ctx, status, OrderStatus)
	if err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.FindByStatusID(ctx, orderStatus.ID)
	if err != nil {
		if errors.Is(err, &core.NotFoundError{}) {
			return nil, core.NewLogicalError(err, orderServiceCode, err.Error())
		}
		return nil, core.NewTechnicalError(err, orderServiceCode, "Ошибка при поиске заказа у клиента")
	}
	return orders, nil
}

func (s *orderService) GetByManagerID(ctx context.Context, managerID uint) ([]Order, error) {
	orders, err := s.orderRepo.FindByManagerID(ctx, managerID)
	if err != nil {
		if errors.Is(err, &core.NotFoundError{}) {
			return nil, core.NewLogicalError(err, orderServiceCode, err.Error())
		}
		return nil, core.NewTechnicalError(err, orderServiceCode, "Ошибка при поиске заказа у клиента")
	}
	return orders, nil
}

func (s *orderService) GetByManagerIDAndStatus(ctx context.Context, managerID uint, status string) ([]Order, error) {
	orderStatus, err := s.enumValueService.GetByCodeAndEnumCode(ctx, status, OrderStatus)
	if err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.FindByManagerIDAndStatusID(ctx, managerID, orderStatus.ID)
	if err != nil {
		if errors.Is(err, &core.NotFoundError{}) {
			return nil, core.NewLogicalError(err, orderServiceCode, err.Error())
		}
		return nil, core.NewTechnicalError(err, orderServiceCode, "Ошибка при поиске заказа у клиента")
	}
	return orders, nil
}
