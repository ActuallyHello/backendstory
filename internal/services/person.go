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
	personServiceCode = "PERSON_SERVICE"
)

type PersonService interface {
	Create(ctx context.Context, person entities.Person) (entities.Person, error)
	Update(ctx context.Context, person entities.Person) (entities.Person, error)
	Delete(ctx context.Context, person entities.Person, soft bool) error

	GetAll(ctx context.Context) ([]entities.Person, error)
	GetById(ctx context.Context, id uint) (entities.Person, error)
	GetByUserLogin(ctx context.Context, userLogin string) (entities.Person, error)
}

type personService struct {
	personRepo repositories.PersonRepository
}

func NewPersonService(
	personRepo repositories.PersonRepository,
) *personService {
	return &personService{
		personRepo: personRepo,
	}
}

// Create создает новую запись Person
func (s *personService) Create(ctx context.Context, person entities.Person) (entities.Person, error) {
	// Проверка существования Person с таким UserLogin
	existingByUserLogin, err := s.personRepo.FindByUserLogin(ctx, person.UserLogin)
	if err != nil && errors.Is(err, &appError.TechnicalError{}) {
		return entities.Person{}, err
	}
	if existingByUserLogin.ID > 0 {
		slog.Error("Person already exists for this user!", "error", err, "user_login", person.UserLogin)
		return entities.Person{}, appError.NewLogicalError(nil, personServiceCode, fmt.Sprintf("Person with user_login = %s already exists!", person.UserLogin))
	}

	// Создаем запись
	created, err := s.personRepo.Create(ctx, person)
	if err != nil {
		slog.Error("Create person failed", "error", err, "user_login", person.UserLogin, "phone", person.Phone)
		return entities.Person{}, appError.NewTechnicalError(err, personServiceCode, err.Error())
	}
	slog.Info("Person created", "firstname", created.Firstname, "lastname", created.Lastname, "user_login", created.UserLogin)
	return created, nil
}

// Update обновляет существующую запись Person
func (s *personService) Update(ctx context.Context, person entities.Person) (entities.Person, error) {
	// Проверяем существование Person
	_, err := s.personRepo.FindById(ctx, person.ID)
	if err != nil {
		return entities.Person{}, err
	}

	updated, err := s.personRepo.Update(ctx, person)
	if err != nil {
		slog.Error("Update person failed", "error", err, "personID", person.ID, "user_login", person.UserLogin)
		return entities.Person{}, appError.NewTechnicalError(err, personServiceCode, err.Error())
	}
	return updated, nil
}

// Delete удаляет Person (мягко или полностью)
func (s *personService) Delete(ctx context.Context, person entities.Person, soft bool) error {
	err := s.personRepo.Delete(ctx, person, soft)
	if err != nil {
		slog.Error("Failed to delete person", "error", err, "personID", person.ID, "soft", soft)
		return appError.NewTechnicalError(err, personServiceCode, err.Error())
	}
	slog.Info("Deleted person", "personID", person.ID, "soft", soft)
	return nil
}

// FindById ищет Person по ID
func (s *personService) GetAll(ctx context.Context) ([]entities.Person, error) {
	persons, err := s.personRepo.FindAll(ctx)
	if err != nil {
		slog.Error("Failed to find persons", "error", err)
		return nil, appError.NewTechnicalError(err, personServiceCode, err.Error())
	}
	return persons, nil
}

// FindById ищет Person по ID
func (s *personService) GetById(ctx context.Context, id uint) (entities.Person, error) {
	person, err := s.personRepo.FindById(ctx, id)
	if err != nil {
		slog.Error("Failed to find person by ID", "error", err, "id", id)
		if errors.Is(err, &common.NotFoundError{}) {
			return entities.Person{}, appError.NewLogicalError(err, personServiceCode, err.Error())
		}
		return entities.Person{}, appError.NewTechnicalError(err, personServiceCode, err.Error())
	}
	return person, nil
}

// FindByUserLogin ищет Person по UserLogin
func (s *personService) GetByUserLogin(ctx context.Context, userLogin string) (entities.Person, error) {
	person, err := s.personRepo.FindByUserLogin(ctx, userLogin)
	if err != nil {
		slog.Error("Failed to find person by user ID", "error", err, "user_login", userLogin)
		if errors.Is(err, &common.NotFoundError{}) {
			return entities.Person{}, appError.NewLogicalError(err, personServiceCode, err.Error())
		}
		return entities.Person{}, appError.NewTechnicalError(err, personServiceCode, err.Error())
	}
	return person, nil
}
