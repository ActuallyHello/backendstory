package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	appError "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/ActuallyHello/backendstory/internal/store/repositories"
	"github.com/ActuallyHello/backendstory/internal/store/repositories/common"
)

const (
	personServiceCode = "PERSON_SERVICE"
)

type PersonService interface {
	BaseService[entities.Person]
	Create(ctx context.Context, person entities.Person) (entities.Person, error)
	Update(ctx context.Context, person entities.Person) (entities.Person, error)
	Delete(ctx context.Context, person entities.Person, soft bool) error

	GetByUserLogin(ctx context.Context, userLogin string) (entities.Person, error)
}

type personService struct {
	BaseServiceImpl[entities.Person]
	personRepo repositories.PersonRepository
}

func NewPersonService(
	personRepo repositories.PersonRepository,
) *personService {
	return &personService{
		BaseServiceImpl: *NewBaseServiceImpl(personRepo),
		personRepo:      personRepo,
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
	_, err := s.personRepo.FindByID(ctx, person.ID)
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
	if soft {
		person.DeletedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		if _, err := s.Update(ctx, person); err != nil {
			slog.Error("Failed to soft delete person", "error", err, "personID", person.ID, "soft", soft)
			return err
		}
	} else {
		err := s.personRepo.Delete(ctx, person)
		if err != nil {
			slog.Error("Failed to delete person", "error", err, "personID", person.ID, "soft", soft)
			return appError.NewTechnicalError(err, personServiceCode, err.Error())
		}
	}
	slog.Info("Deleted person", "personID", person.ID, "soft", soft)
	return nil
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
