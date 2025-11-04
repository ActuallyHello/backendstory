package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	appError "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
)

const (
	roleServiceCode = "ROLE_SERVICE"
)

type RoleService interface {
	Create(ctx context.Context, role entities.Role) (entities.Role, error)
	Update(ctx context.Context, role entities.Role) (entities.Role, error)
	Delete(ctx context.Context, role entities.Role) error

	GetAll(ctx context.Context) ([]entities.Role, error)
	GetById(ctx context.Context, id uint) (entities.Role, error)
	GetByCode(ctx context.Context, code string) (entities.Role, error)
}

type roleService struct {
	roleRepo repositories.RoleRepository
}

func NewRoleService(
	roleRepo repositories.RoleRepository,
) *roleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}

func (s *roleService) Create(ctx context.Context, role entities.Role) (entities.Role, error) {
	// Проверка существования с таким кодом
	existing, err := s.GetByCode(ctx, role.Code)
	if err != nil && errors.Is(err, &appError.TechnicalError{}) {
		return entities.Role{}, err
	}
	if existing.ID > 0 {
		slog.Error("Role already exists!", "error", err, "code", role.Code)
		return entities.Role{}, appError.NewLogicalError(nil, roleServiceCode, fmt.Sprintf("Role with code = %s already exists!", role.Code))
	}

	// Создаем запись
	created, err := s.roleRepo.Create(ctx, role)
	if err != nil {
		slog.Error("Create role failed", "error", err, "code", role.Code)
		return entities.Role{}, appError.NewTechnicalError(err, roleServiceCode, err.Error())
	}
	slog.Info("Role created", "code", created.Code)
	return created, nil
}

func (s *roleService) Update(ctx context.Context, role entities.Role) (entities.Role, error) {
	existing, err := s.GetById(ctx, role.ID)
	if err != nil {
		return entities.Role{}, err
	}

	if existing.Code != role.Code {
		existingByCode, err := s.GetByCode(ctx, role.Code)
		if err != nil && errors.Is(err, &appError.TechnicalError{}) {
			return entities.Role{}, err
		}
		if existingByCode.ID > 0 {
			slog.Error("Role already exists!", "error", err, "code", role.Code)
			return entities.Role{}, appError.NewLogicalError(err, roleServiceCode, fmt.Sprintf("Role with code = %s already exists!", role.Code))
		}
	}

	if role.Code != "" {
		existing.Code = role.Code
	}
	if role.Label != "" {
		existing.Label = role.Label
	}

	updated, err := s.roleRepo.Update(ctx, existing)
	if err != nil {
		slog.Error("Update role failed", "error", err, "code", role.Code)
		return entities.Role{}, err
	}
	return updated, nil
}

func (s *roleService) Delete(ctx context.Context, role entities.Role) error {
	err := s.roleRepo.Delete(ctx, role)
	if err != nil {
		slog.Error("Failed to delete role", "error", err, "id", role.ID)
		return appError.NewTechnicalError(err, roleServiceCode, err.Error())
	}
	slog.Info("Deleted role", "code", role.Code)
	return nil
}

func (s *roleService) GetAll(ctx context.Context) ([]entities.Role, error) {
	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to find roles", "error", err)
		return nil, appError.NewTechnicalError(err, roleServiceCode, err.Error())
	}
	return roles, nil
}

func (s *roleService) GetById(ctx context.Context, id uint) (entities.Role, error) {
	role, err := s.roleRepo.FindById(ctx, id)
	if err != nil {
		slog.Error("Failed to find role by ID", "error", err, "id", id)
		if errors.Is(err, &common.NotFoundError{}) {
			return entities.Role{}, appError.NewLogicalError(err, roleServiceCode, err.Error())
		}
		return entities.Role{}, appError.NewTechnicalError(err, roleServiceCode, err.Error())
	}
	return role, nil
}

func (s *roleService) GetByCode(ctx context.Context, code string) (entities.Role, error) {
	role, err := s.roleRepo.FindByCode(ctx, code)
	if err != nil {
		slog.Error("Failed to find role by code", "error", err, "code", code)
		if errors.Is(err, &common.NotFoundError{}) {
			return entities.Role{}, appError.NewLogicalError(err, roleServiceCode, err.Error())
		}
		return entities.Role{}, appError.NewTechnicalError(err, roleServiceCode, err.Error())
	}
	return role, nil
}
