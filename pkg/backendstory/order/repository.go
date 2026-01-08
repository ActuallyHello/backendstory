package order

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/gorm"
)

type OrderRepository interface {
	core.BaseRepository[Order]

	FindByStatusID(ctx context.Context, statusID uint) ([]Order, error)
	FindByClientID(ctx context.Context, clientID uint) ([]Order, error)
	FindByManagerID(ctx context.Context, managerID uint) ([]Order, error)
	FindByManagerIDAndStatusID(ctx context.Context, managerID, statusID uint) ([]Order, error)
}

type orderRepository struct {
	core.BaseRepositoryImpl[Order]
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{
		BaseRepositoryImpl: *core.NewBaseRepositoryImpl[Order](db),
	}
}

func (r *orderRepository) FindByStatusID(ctx context.Context, statusID uint) ([]Order, error) {
	var orders []Order
	if err := r.GetDB(ctx).Where("STATUSID = ?", statusID).Find(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.NewNotFoundError("Не существует заказов с данным статусом")
		}
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) FindByClientID(ctx context.Context, clientID uint) ([]Order, error) {
	var orders []Order
	if err := r.GetDB(ctx).Where("CLIENTID = ?", clientID).Find(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.NewNotFoundError("Не существует заказов у данного клиента")
		}
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) FindByManagerID(ctx context.Context, managerID uint) ([]Order, error) {
	var orders []Order
	if err := r.GetDB(ctx).Where("MANAGERID = ?", managerID).Find(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.NewNotFoundError("Не существует заказов у данного менеджера")
		}
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) FindByManagerIDAndStatusID(ctx context.Context, managerID, statusID uint) ([]Order, error) {
	var orders []Order
	if err := r.GetDB(ctx).Where("MANAGERID = ? AND STATUSID = ?", managerID, statusID).Find(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.NewNotFoundError("Не существует заказов у данного менеджера")
		}
		return nil, err
	}
	return orders, nil
}
