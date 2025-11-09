package container

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/ActuallyHello/backendstory/internal/config"
	"github.com/ActuallyHello/backendstory/internal/server/handlers"
	"github.com/ActuallyHello/backendstory/internal/services"
	"github.com/ActuallyHello/backendstory/internal/services/auth"
	"github.com/ActuallyHello/backendstory/internal/store/repositories"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AppContainer struct {
	// application
	ctx context.Context

	// database
	db *gorm.DB

	// repositoriest
	enumRepo      repositories.EnumRepository
	enumValueRepo repositories.EnumValueRepository
	personRepo    repositories.PersonRepository

	// services
	enumService      services.EnumService
	enumValueService services.EnumValueService
	personService    services.PersonService

	// handlers
	authHandler      *handlers.AuthHandler
	enumHandler      *handlers.EnumHandler
	enumValueHandler *handlers.EnumValueHandler
	personHandler    *handlers.PersonHandler

	// auth
	authService auth.AuthService
}

func NewAppContainer(appConfig *config.ApplicationConfig) *AppContainer {
	// application
	ctx := context.Background()

	// database
	dsn := constructDSN(appConfig.DatabaseConfig)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Error while establish database connection", "err", err)
		log.Fatal(err)
	}

	// repositories
	enumRepo := repositories.NewEnumRepository(db)
	enumValueRepo := repositories.NewEnumValueRepository(db)
	personRepo := repositories.NewPersonRepository(db)

	// services
	enumService := services.NewEnumService(enumRepo)
	enumValueService := services.NewEnumValueService(enumValueRepo)
	personService := services.NewPersonService(personRepo)

	// auth
	keycloakService, err := auth.NewKeycloakService(ctx, appConfig.KeycloakConfig)
	if err != nil {
		slog.Error("Error while creating keycloak connection", "err", err)
		log.Fatal(err)
	}

	// handlers
	enumHandler := handlers.NewEnumHandler(enumService)
	enumValueHandler := handlers.NewEnumValueHandler(enumValueService)
	personHandler := handlers.NewPersonHandler(personService)
	authHandler := handlers.NewAuthHandler(keycloakService)

	return &AppContainer{
		// application
		ctx: ctx,

		// database
		db: db,

		// repositoriest
		enumRepo:      enumRepo,
		enumValueRepo: enumValueRepo,
		personRepo:    personRepo,

		// services
		enumService:      enumService,
		enumValueService: enumValueService,
		personService:    personService,

		// handlers
		authHandler:      authHandler,
		enumHandler:      enumHandler,
		enumValueHandler: enumValueHandler,
		personHandler:    personHandler,

		// auth
		authService: keycloakService,
	}
}

func constructDSN(databaseConfig *config.DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		databaseConfig.Username,
		databaseConfig.Password,
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.Database,
	)
}

// Application
func (c *AppContainer) GetContext() context.Context {
	return c.ctx
}

// Database
func (c *AppContainer) GetDB() *gorm.DB {
	return c.db
}

// Repositories
func (c *AppContainer) GetEnumRepository() repositories.EnumRepository {
	return c.enumRepo
}

func (c *AppContainer) GetEnumValueRepository() repositories.EnumValueRepository {
	return c.enumValueRepo
}

func (c *AppContainer) GetPersonRepository() repositories.PersonRepository {
	return c.personRepo
}

// Services
func (c *AppContainer) GetEnumService() services.EnumService {
	return c.enumService
}

func (c *AppContainer) GetEnumValueService() services.EnumValueService {
	return c.enumValueService
}

func (c *AppContainer) GetPersonService() services.PersonService {
	return c.personService
}

// Handlers
func (c *AppContainer) GetAuthHandler() *handlers.AuthHandler {
	return c.authHandler
}

func (c *AppContainer) GetEnumHandler() *handlers.EnumHandler {
	return c.enumHandler
}

func (c *AppContainer) GetEnumValueHandler() *handlers.EnumValueHandler {
	return c.enumValueHandler
}

func (c *AppContainer) GetPersonHandler() *handlers.PersonHandler {
	return c.personHandler
}

// Auth
func (c *AppContainer) GetAuthService() auth.AuthService {
	return c.authService
}

// Close освобождает ресурсы
func (c *AppContainer) Close() error {
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// HealthCheck проверяет доступность базы данных
func (c *AppContainer) HealthCheck() error {
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	}
	return fmt.Errorf("database connection is not initialized")
}
