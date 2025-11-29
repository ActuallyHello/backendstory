package container

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/auth"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/cart"
	cartitem "github.com/ActuallyHello/backendstory/pkg/backendstory/cart_item"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/category"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enum"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/person"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/product"
	productmedia "github.com/ActuallyHello/backendstory/pkg/backendstory/product_media"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/purchase"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/resources"
	"github.com/ActuallyHello/backendstory/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AppContainer struct {
	// application
	ctx context.Context

	// database
	db *gorm.DB

	// repositoriest
	enumRepo         enum.EnumRepository
	enumValueRepo    enumvalue.EnumValueRepository
	personRepo       person.PersonRepository
	categoryRepo     category.CategoryRepository
	productRepo      product.ProductRepository
	productMediaRepo productmedia.ProductMediaRepository
	cartRepo         cart.CartRepository
	cartItemRepo     cartitem.CartItemRepository

	// services
	enumService         enum.EnumService
	enumValueService    enumvalue.EnumValueService
	personService       person.PersonService
	categoryService     category.CategoryService
	productService      product.ProductService
	productMediaService productmedia.ProductMediaService
	cartService         cart.CartService
	cartItemService     cartitem.CartItemService
	purchaseService     purchase.PurchaseService

	// resources
	fileService resources.FileService

	// handlers
	authHandler         *auth.AuthHandler
	enumHandler         *enum.EnumHandler
	enumValueHandler    *enumvalue.EnumValueHandler
	personHandler       *person.PersonHandler
	categoryHandler     *category.CategoryHandler
	productHandler      *product.ProductHandler
	productMediaHandler *productmedia.ProductMediaHandler
	cartHandler         *cart.CartHandler
	cartItemHandler     *cartitem.CartItemHandler
	purchaseHandler     *purchase.PurchaseHandler

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
	enumRepo := enum.NewEnumRepository(db)
	enumValueRepo := enumvalue.NewEnumValueRepository(db)
	personRepo := person.NewPersonRepository(db)
	categoryRepo := category.NewCategoryRepository(db)
	productRepo := product.NewProductRepository(db)
	productMediaRepo := productmedia.NewProductMediaRepository(db)
	cartRepo := cart.NewCartRepository(db)
	cartItemRepo := cartitem.NewCartItemRepository(db)

	// services
	enumService := enum.NewEnumService(enumRepo)
	enumValueService := enumvalue.NewEnumValueService(enumValueRepo)
	personService := person.NewPersonService(personRepo)
	categoryService := category.NewCategoryService(categoryRepo)
	productService := product.NewProductService(productRepo, enumService, enumValueService)
	productMediaService := productmedia.NewProductMediaService(productMediaRepo)
	cartServices := cart.NewCartService(cartRepo)
	cartItemServices := cartitem.NewCartItemService(cartItemRepo, enumService, enumValueService)
	purchaseService := purchase.NewPurchaseService(cartServices, cartItemServices, productService, enumService, enumValueService)

	// auth
	keycloakService, err := auth.NewKeycloakService(ctx, appConfig.KeycloakConfig)
	if err != nil {
		slog.Error("Error while creating keycloak connection", "err", err)
		log.Fatal(err)
	}

	// handlers
	enumHandler := enum.NewEnumHandler(enumService)
	enumValueHandler := enumvalue.NewEnumValueHandler(enumValueService)
	personHandler := person.NewPersonHandler(personService)
	authHandler := auth.NewAuthHandler(keycloakService)
	categoryHandler := category.NewCategoryHandler(categoryService)
	productHandler := product.NewProductHandler(productService, enumValueService, categoryService)
	cartHandler := cart.NewCartHandler(cartServices)
	cartItemHandler := cartitem.NewCartItemHandler(cartItemServices)
	purchaseHandler := purchase.NewPurchaseHandler(purchaseService)

	// refactor
	productMediaHandler := productmedia.NewProductMediaHandler(productMediaService, productService, resources.FileService{}, "static/media")

	return &AppContainer{
		// application
		ctx: ctx,

		// database
		db: db,

		// repositoriest
		enumRepo:         enumRepo,
		enumValueRepo:    enumValueRepo,
		personRepo:       personRepo,
		categoryRepo:     categoryRepo,
		productRepo:      productRepo,
		productMediaRepo: productMediaRepo,
		cartRepo:         cartRepo,
		cartItemRepo:     cartItemRepo,

		// services
		enumService:         enumService,
		enumValueService:    enumValueService,
		personService:       personService,
		categoryService:     categoryService,
		productService:      productService,
		productMediaService: productMediaService,
		cartService:         cartServices,
		cartItemService:     cartItemServices,
		purchaseService:     purchaseService,

		fileService: resources.FileService{},

		// handlers
		authHandler:         authHandler,
		enumHandler:         enumHandler,
		enumValueHandler:    enumValueHandler,
		personHandler:       personHandler,
		categoryHandler:     categoryHandler,
		productHandler:      productHandler,
		productMediaHandler: productMediaHandler,
		cartHandler:         cartHandler,
		cartItemHandler:     cartItemHandler,
		purchaseHandler:     purchaseHandler,

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
func (c *AppContainer) GetEnumRepository() enum.EnumRepository {
	return c.enumRepo
}

func (c *AppContainer) GetEnumValueRepository() enumvalue.EnumValueRepository {
	return c.enumValueRepo
}

func (c *AppContainer) GetPersonRepository() person.PersonRepository {
	return c.personRepo
}

func (c *AppContainer) GetCategoryRepository() category.CategoryRepository {
	return c.categoryRepo
}

func (c *AppContainer) GetProductRepository() product.ProductRepository {
	return c.productRepo
}

func (c *AppContainer) GetProductMediaRepository() productmedia.ProductMediaRepository {
	return c.productMediaRepo
}

func (c *AppContainer) GetCartRepository() cart.CartRepository {
	return c.cartRepo
}

func (c *AppContainer) GetCartItemRepository() cartitem.CartItemRepository {
	return c.cartItemRepo
}

// Services
func (c *AppContainer) GetEnumService() enum.EnumService {
	return c.enumService
}

func (c *AppContainer) GetEnumValueService() enumvalue.EnumValueService {
	return c.enumValueService
}

func (c *AppContainer) GetPersonService() person.PersonService {
	return c.personService
}

func (c *AppContainer) GetCategoryService() category.CategoryService {
	return c.categoryService
}

func (c *AppContainer) GetProductService() product.ProductService {
	return c.productService
}

func (c *AppContainer) GetProductMediaService() productmedia.ProductMediaService {
	return c.productMediaService
}

func (c *AppContainer) GetCartService() cart.CartService {
	return c.cartService
}

func (c *AppContainer) GetCartItemService() cartitem.CartItemService {
	return c.cartItemService
}

func (c *AppContainer) GetPurchaseService() purchase.PurchaseService {
	return c.purchaseService
}

func (c *AppContainer) GetFileService() resources.FileService {
	return c.fileService
}

// Handlers
func (c *AppContainer) GetAuthHandler() *auth.AuthHandler {
	return c.authHandler
}

func (c *AppContainer) GetEnumHandler() *enum.EnumHandler {
	return c.enumHandler
}

func (c *AppContainer) GetEnumValueHandler() *enumvalue.EnumValueHandler {
	return c.enumValueHandler
}

func (c *AppContainer) GetPersonHandler() *person.PersonHandler {
	return c.personHandler
}

func (c *AppContainer) GetCategoryHandler() *category.CategoryHandler {
	return c.categoryHandler
}

func (c *AppContainer) GetProductHandler() *product.ProductHandler {
	return c.productHandler
}

func (c *AppContainer) GetProductMediaHandler() *productmedia.ProductMediaHandler {
	return c.productMediaHandler
}

func (c *AppContainer) GetCartHandler() *cart.CartHandler {
	return c.cartHandler
}

func (c *AppContainer) GetCartItemHandler() *cartitem.CartItemHandler {
	return c.cartItemHandler
}

func (c *AppContainer) GetPurchaseHandler() *purchase.PurchaseHandler {
	return c.purchaseHandler
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
