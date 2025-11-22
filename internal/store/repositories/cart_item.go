package repositories

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	BaseRepository[entities.CartItem]

	FindByCartID(ctx context.Context, cartID uint) ([]entities.CartItem, error)
}

type cartItemRepository struct {
	BaseRepositoryImpl[entities.CartItem]
}

func NewCartItemRepository(db *gorm.DB) *cartItemRepository {
	return &cartItemRepository{
		BaseRepositoryImpl: *NewBaseRepositoryImpl[entities.CartItem](db),
	}
}

// FindByCartID ищет CartItem по Cart
func (r *cartItemRepository) FindByCartID(ctx context.Context, cartID uint) ([]entities.CartItem, error) {
	var cartItems []entities.CartItem
	if err := r.db.WithContext(ctx).Where("CARTID = ?", cartID).Find(&cartItems).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewNotFoundError("CartItem not found by code")
		}
		return nil, err
	}
	return cartItems, nil
}
