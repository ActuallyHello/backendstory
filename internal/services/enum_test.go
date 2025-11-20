package services

import (
	"context"
	"testing"

	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEnumRepository struct {
	mock.Mock
}

func (m *MockEnumRepository) Create(ctx context.Context, enum entities.Enum) (entities.Enum, error) {
	args := m.Called(ctx, enum)
	return args.Get(0).(entities.Enum), args.Error(1)
}

func (m *MockEnumRepository) Update(ctx context.Context, enum entities.Enum) (entities.Enum, error) {
	args := m.Called(ctx, enum)
	return args.Get(0).(entities.Enum), args.Error(1)
}

func (m *MockEnumRepository) Delete(ctx context.Context, enum entities.Enum) error {
	args := m.Called(ctx, enum)
	return args.Error(0)
}

func (m *MockEnumRepository) FindByID(ctx context.Context, id uint) (entities.Enum, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.Enum), args.Error(1)
}

func (m *MockEnumRepository) FindAll(ctx context.Context) ([]entities.Enum, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Enum), args.Error(1)
}

func (m *MockEnumRepository) FindByCode(ctx context.Context, code string) (entities.Enum, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(entities.Enum), args.Error(1)
}

func (m *MockEnumRepository) FindWithSearchCriteria(ctx context.Context, criteria dto.SearchCriteria) ([]entities.Enum, error) {
	args := m.Called(ctx, criteria)
	return args.Get(0).([]entities.Enum), args.Error(1)
}

func (m *MockEnumRepository) Count(ctx context.Context, criteria *dto.SearchCriteria) (int64, error) {
	args := m.Called(ctx, criteria)
	return args.Get(0).(int64), args.Error(1)
}

func TestEnumService_Create_Success(t *testing.T) {
	mockRepo := new(MockEnumRepository)
	service := NewEnumService(mockRepo)

	ctx := context.Background()
	newEnum := entities.Enum{
		Code:  "TEST_CODE",
		Label: "Test label",
	}

	mockRepo.On("FindByCode", ctx, "TEST_CODE").Return(entities.Enum{}, common.NewNotFoundError("enum not found"))

	createdEnum := entities.Enum{
		Base: entities.Base{
			ID: 1,
		},
		Code:  "TEST_CODE",
		Label: "Test label",
	}
	mockRepo.On("Create", ctx, newEnum).Return(createdEnum, nil)

	result, err := service.Create(ctx, newEnum)

	assert.NoError(t, err)
	assert.Equal(t, createdEnum, result)
	mockRepo.AssertExpectations(t)

}
