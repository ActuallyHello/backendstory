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
	"github.com/ActuallyHello/backendstory/pkg/backendstory/order"
	orderitem "github.com/ActuallyHello/backendstory/pkg/backendstory/order_item"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/person"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/product"
	productmedia "github.com/ActuallyHello/backendstory/pkg/backendstory/product_media"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/resources"
	"github.com/ActuallyHello/backendstory/pkg/config"
	"github.com/ActuallyHello/backendstory/pkg/core"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AppContainer struct {
	// application
	ctx    context.Context
	cancel context.CancelFunc

	// database
	db        *gorm.DB
	txManager core.TxManager

	// repositoriest
	enumRepo         enum.EnumRepository
	enumValueRepo    enumvalue.EnumValueRepository
	personRepo       person.PersonRepository
	categoryRepo     category.CategoryRepository
	productRepo      product.ProductRepository
	productMediaRepo productmedia.ProductMediaRepository
	cartRepo         cart.CartRepository
	cartItemRepo     cartitem.CartItemRepository
	orderRepo        order.OrderRepository
	orderItemRepo    orderitem.OrderItemRepository

	// services
	enumService         enum.EnumService
	enumValueService    enumvalue.EnumValueService
	personService       person.PersonService
	categoryService     category.CategoryService
	productService      product.ProductService
	productMediaService productmedia.ProductMediaService
	cartService         cart.CartService
	cartItemService     cartitem.CartItemService
	orderItemService    orderitem.OrderItemService
	orderService        order.OrderService

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
	orderHandler        *order.OrderHandler
	orderItemHandler    *orderitem.OrderItemHandler

	// auth
	authService auth.AuthService
}

func NewAppContainer(ctx context.Context, appConfig *config.ApplicationConfig) (*AppContainer, error) {
	// application
	appCtx, cancel := context.WithCancel(ctx)

	// database
	dsn := constructDSN(appConfig.DatabaseConfig)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Error while establish database connection", "err", err)
		log.Fatal(err)
	}
	txManager := core.NewGormTxManager(db)

	// repositories
	enumRepo := enum.NewEnumRepository(db)
	enumValueRepo := enumvalue.NewEnumValueRepository(db)
	personRepo := person.NewPersonRepository(db)
	categoryRepo := category.NewCategoryRepository(db)
	productRepo := product.NewProductRepository(db)
	productMediaRepo := productmedia.NewProductMediaRepository(db)
	cartRepo := cart.NewCartRepository(db)
	cartItemRepo := cartitem.NewCartItemRepository(db)
	orderRepo := order.NewOrderRepository(db)
	orderItemRepo := orderitem.NewOrderItemRepository(db)

	// services
	enumService := enum.NewEnumService(enumRepo)
	enumValueService := enumvalue.NewEnumValueService(enumValueRepo, enumService)
	personService := person.NewPersonService(personRepo)
	categoryService := category.NewCategoryService(categoryRepo)
	productService := product.NewProductService(productRepo, enumService, enumValueService)
	productMediaService := productmedia.NewProductMediaService(productMediaRepo)
	cartServices := cart.NewCartService(cartRepo)
	cartItemService := cartitem.NewCartItemService(cartItemRepo, enumService, enumValueService, productService)
	orderItemService := orderitem.NewOrderItemService(orderItemRepo, txManager, enumService, enumValueService, productService, cartItemService)
	orderService := order.NewOrderService(orderRepo, txManager, enumService, enumValueService, orderItemService)

	// auth
	keycloakService, err := auth.NewKeycloakService(appCtx, appConfig.KeycloakConfig)
	if err != nil {
		slog.Error("Error while creating keycloak connection", "err", err)
		log.Fatal(err)
	}

	// resources
	fileService := resources.NewFileService()

	// handlers
	enumHandler := enum.NewEnumHandler(enumService)
	enumValueHandler := enumvalue.NewEnumValueHandler(enumValueService)
	personHandler := person.NewPersonHandler(personService)
	authHandler := auth.NewAuthHandler(keycloakService)
	categoryHandler := category.NewCategoryHandler(categoryService)
	productHandler := product.NewProductHandler(productService, enumValueService, categoryService)
	cartHandler := cart.NewCartHandler(cartServices)
	cartItemHandler := cartitem.NewCartItemHandler(cartItemService)
	orderItemHandler := orderitem.NewOrderItemHandler(orderItemService, enumValueService)
	orderHandler := order.NewOrderHandler(orderService, personService, enumValueService)
	productMediaHandler := productmedia.NewProductMediaHandler(productMediaService, productService, fileService, appConfig.ServerConfig.StaticFilesPath)

	return &AppContainer{
		// application
		ctx:    appCtx,
		cancel: cancel,

		// database
		db:        db,
		txManager: txManager,

		// repositoriest
		enumRepo:         enumRepo,
		enumValueRepo:    enumValueRepo,
		personRepo:       personRepo,
		categoryRepo:     categoryRepo,
		productRepo:      productRepo,
		productMediaRepo: productMediaRepo,
		cartRepo:         cartRepo,
		cartItemRepo:     cartItemRepo,
		orderRepo:        orderRepo,
		orderItemRepo:    orderItemRepo,

		// services
		enumService:         enumService,
		enumValueService:    enumValueService,
		personService:       personService,
		categoryService:     categoryService,
		productService:      productService,
		productMediaService: productMediaService,
		cartService:         cartServices,
		cartItemService:     cartItemService,
		orderItemService:    orderItemService,
		orderService:        orderService,

		// resources
		fileService: fileService,

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
		orderHandler:        orderHandler,
		orderItemHandler:    orderItemHandler,

		// auth
		authService: keycloakService,
	}, nil
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

// Close освобождает ресурсы
func (c *AppContainer) Close() {
	slog.Info("Closing application resources")

	c.cancel()

	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			slog.Error("failed attempt to close database connection", "error", err)
		}
		if err := sqlDB.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}

	slog.Info("All resources closed")
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
	return core.NewTechnicalError(nil, "CONTAINER", "database connection is not initialized")
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

func (c *AppContainer) GetOrderRepository() order.OrderRepository {
	return c.orderRepo
}

func (c *AppContainer) GetOrderItemRepository() orderitem.OrderItemRepository {
	return c.orderItemRepo
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

func (c *AppContainer) GetFileService() resources.FileService {
	return c.fileService
}

func (c *AppContainer) GetOrderService() order.OrderService {
	return c.orderService
}

func (c *AppContainer) GetOrderItemService() orderitem.OrderItemService {
	return c.orderItemService
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

func (c *AppContainer) GetOrderHandler() *order.OrderHandler {
	return c.orderHandler
}

func (c *AppContainer) GetOrderItemHandler() *orderitem.OrderItemHandler {
	return c.orderItemHandler
}

// Auth
func (c *AppContainer) GetAuthService() auth.AuthService {
	return c.authService
}
