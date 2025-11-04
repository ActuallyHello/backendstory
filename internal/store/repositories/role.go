package repositories

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(ctx context.Context, role entities.Role) (entities.Role, error)
	Update(ctx context.Context, role entities.Role) (entities.Role, error)
	Delete(ctx context.Context, role entities.Role) error

	FindAll(ctx context.Context) ([]entities.Role, error)
	FindById(ctx context.Context, id uint) (entities.Role, error)
	FindByCode(ctx context.Context, code string) (entities.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role entities.Role) (entities.Role, error) {
	if err := r.db.WithContext(ctx).Create(&role).Error; err != nil {
		return entities.Role{}, err
	}
	return role, nil
}

func (r *roleRepository) Delete(ctx context.Context, role entities.Role) error {
	if err := r.db.WithContext(ctx).Delete(&role).Error; err != nil {
		return err
	}
	return nil
}

func (r *roleRepository) Update(ctx context.Context, role entities.Role) (entities.Role, error) {
	if err := r.db.WithContext(ctx).Save(&role).Error; err != nil {
		return entities.Role{}, err
	}
	return role, nil
}

func (r *roleRepository) FindAll(ctx context.Context) ([]entities.Role, error) {
	var roles []entities.Role
	if err := r.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) FindById(ctx context.Context, id uint) (entities.Role, error) {
	var role entities.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Role{}, common.NewNotFoundError("role not found by id")
		}
		return entities.Role{}, err
	}
	return role, nil
}

func (r *roleRepository) FindByCode(ctx context.Context, code string) (entities.Role, error) {
	var role entities.Role
	if err := r.db.WithContext(ctx).Where("CODE = ?", code).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Role{}, common.NewNotFoundError("role not found by code")
		}
		return entities.Role{}, err
	}
	return role, nil
}
