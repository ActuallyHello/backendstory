package enum

import (
	"context"
	"errors"

	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	enumServiceCode = "ENUMERATION_SERVICE"
)

type EnumService interface {
	core.BaseService[Enum]

	Create(ctx context.Context, enum Enum) (Enum, error)
	Update(ctx context.Context, enum Enum) (Enum, error)
	Delete(ctx context.Context, enum Enum) error

	GetByCode(ctx context.Context, code string) (Enum, error)
}

type enumService struct {
	core.BaseServiceImpl[Enum]
	enumRepo EnumRepository
}

func NewEnumService(
	enumRepo EnumRepository,
) *enumService {
	return &enumService{
		BaseServiceImpl: *core.NewBaseServiceImpl(enumRepo),
		enumRepo:        enumRepo,
	}
}

func (s *enumService) Create(ctx context.Context, enum Enum) (Enum, error) {
	existing, err := s.GetByCode(ctx, enum.Code)
	if err != nil && errors.Is(err, &core.TechnicalError{}) {
		return Enum{}, err
	}
	if existing.ID > 0 {
		return Enum{}, core.NewLogicalError(nil, enumServiceCode, "Перечисление уже существует")
	}

	created, err := s.GetRepo().Create(ctx, enum)
	if err != nil {
		return Enum{}, core.NewTechnicalError(err, enumServiceCode, "Ошибка при создании перечисления")
	}
	return created, nil
}

func (s *enumService) Update(ctx context.Context, enum Enum) (Enum, error) {
	existing, err := s.GetRepo().FindByID(ctx, enum.ID)
	if err != nil {
		return Enum{}, err
	}

	updated, err := s.GetRepo().Update(ctx, existing)
	if err != nil {
		return Enum{}, core.NewTechnicalError(err, enumServiceCode, "Ошибка при обновлении перечисления")
	}
	return updated, nil
}

func (s *enumService) Delete(ctx context.Context, enum Enum) error {
	err := s.GetRepo().Delete(ctx, enum)
	if err != nil {
		return core.NewTechnicalError(err, enumServiceCode, "Ошибка при удалении перечисления")
	}
	return nil
}

func (s *enumService) GetByCode(ctx context.Context, code string) (Enum, error) {
	enum, err := s.enumRepo.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, &core.NotFoundError{}) {
			return Enum{}, core.NewLogicalError(err, enumServiceCode, err.Error())
		}
		return Enum{}, core.NewTechnicalError(err, enumServiceCode, "Ошибка при получении перечисления по коду")
	}
	return enum, nil
}
