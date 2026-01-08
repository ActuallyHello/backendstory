package cartitem

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	core.BaseRepository[CartItem]

	FindByCartID(ctx context.Context, cartID uint) ([]CartItem, error)
	FindByCartIDAndProductID(ctx context.Context, cartID, productID uint) (CartItem, error)
}

type cartItemRepository struct {
	core.BaseRepositoryImpl[CartItem]
}

func NewCartItemRepository(db *gorm.DB) *cartItemRepository {
	return &cartItemRepository{
		BaseRepositoryImpl: *core.NewBaseRepositoryImpl[CartItem](db),
	}
}

// FindByCartID ищет CartItem по Cart
func (r *cartItemRepository) FindByCartID(ctx context.Context, cartID uint) ([]CartItem, error) {
	var cartItems []CartItem
	if err := r.GetDB(ctx).Where("CARTID = ?", cartID).Find(&cartItems).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.NewNotFoundError("Элементов корзины не найдено")
		}
		return nil, err
	}
	return cartItems, nil
}

// FindByCartID ищет CartItem по Cart
func (r *cartItemRepository) FindByCartIDAndProductID(ctx context.Context, cartID, productID uint) (CartItem, error) {
	var cartItem CartItem
	if err := r.GetDB(ctx).Where("CARTID = ? AND PRODUCTID = ?", cartID, productID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CartItem{}, core.NewNotFoundError("Элементов корзины не найдено")
		}
		return CartItem{}, err
	}
	return cartItem, nil
}
