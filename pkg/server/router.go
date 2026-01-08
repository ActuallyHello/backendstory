package server

import (
	"net/http"
	"os"
	"path/filepath"

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
	"github.com/ActuallyHello/backendstory/pkg/container"
	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	apiV1    = "/api/v1/"
	register = "/register"
	login    = "/login"
	byId     = "/{id}"
)

func SetupRouter(container *container.AppContainer, staticFilesPath string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// custom ErrorHandler
	r.Use(core.ErrorHandler)
	// static files
	setupStaticRoutes(r, staticFilesPath)

	r.Route(apiV1, func(r chi.Router) {

		r.Post(register, container.GetAuthHandler().Register)
		r.Post(login, container.GetAuthHandler().Login)

		// TODO: MAIL FOR ORDER, MAIL FOR APPROVE, RABBITMQ, STATUS MODEL
		// TODO: TEST SCENARIOUS
		// TODO: purchase routes
		// TODO: Order routes, scenarious
		// TODO: enum service by code and by enum value code
		// TODO: order decomposite methods
		// TODO: order item routes

		// TODO: ENUM CONSTAT BY EACH ENTITY.go
		// TODO: DELETE CART ITEM STATUS
		// TODO: common method with approve/cancel order actions

		registerAuthRoutes(r, container.GetAuthService(), container.GetAuthHandler())

		registerEnumRoutes(r, container.GetAuthService(), container.GetEnumHandler())
		registerEnumValuesRoutes(r, container.GetAuthService(), container.GetEnumValueHandler())
		registerPersonRoutes(r, container.GetAuthService(), container.GetPersonHandler())
		registerCategoryRoutes(r, container.GetAuthService(), container.GetCategoryHandler())
		registerProductRoutes(r, container.GetAuthService(), container.GetProductHandler())
		registerProductMediaRoutes(r, container.GetAuthService(), container.GetProductMediaHandler())
		registerCartRoutes(r, container.GetAuthService(), container.GetCartHandler())
		registerCartItemRoutes(r, container.GetAuthService(), container.GetCartItemHandler())
		registerOrderRoutes(r, container.GetAuthService(), container.GetOrderHandler())
		registerOrderItemRoutes(r, container.GetAuthService(), container.GetOrderItemHandler())
	})

	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}

func registerEnumRoutes(r chi.Router, authService auth.AuthService, enumHandler *enum.EnumHandler) {
	r.Route("/enumerations", func(r chi.Router) {
		r.Use(AuthMiddleware(authService, "admin"))

		r.Get("/", enumHandler.GetAll)
		r.Get(byId, enumHandler.GetById)
		r.Get("/code/{code}", enumHandler.GetByCode)
		r.Post("/search", enumHandler.GetWithSearchCriteria)

		// Защищенные маршруты (требуют аутентификации)
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(authService, "admin", "guest"))

			r.Post("/", enumHandler.Create)
			r.Delete(byId, enumHandler.Delete)
		})
	})
}

func registerEnumValuesRoutes(r chi.Router, authService auth.AuthService, enumValueHandler *enumvalue.EnumValueHandler) {
	r.Route("/enumeration-values", func(r chi.Router) {
		r.Use(AuthMiddleware(authService, "admin"))

		r.Get("/", enumValueHandler.GetAll)
		r.Get(byId, enumValueHandler.GetById)
		r.Get("/enumeration/{enumeration_id}", enumValueHandler.GetByEnumId)
		r.Post("/search", enumValueHandler.GetWithSearchCriteria)

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(authService, "admin", "guest"))

			r.Post("/", enumValueHandler.Create)
			r.Delete(byId, enumValueHandler.Delete)
		})
	})
}

func registerPersonRoutes(r chi.Router, authService auth.AuthService, personHandler *person.PersonHandler) {
	r.Route("/persons", func(r chi.Router) {
		r.Use(AuthMiddleware(authService, "admin", "guest"))

		r.Get(byId, personHandler.GetById)
		r.Get("/user/{user_login}", personHandler.GetByUserLogin)
		r.Post("/search", personHandler.GetWithSearchCriteria)

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(authService, "admin"))

			r.Get("/", personHandler.GetAll)
			r.Post("/", personHandler.Create)
			r.Delete(byId, personHandler.Delete)
		})
	})
}

func registerAuthRoutes(r chi.Router, authService auth.AuthService, authHandler *auth.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/token", authHandler.GetHeaderTokenInfo)
		r.Post("/token", authHandler.GetBodyTokenInfo)

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(authService, "admin", "guest"))

			r.Get("/users", authHandler.GetUsers)
			r.Get("/users/{username}", authHandler.GetUser)
			r.Get("/users/{username}/roles", authHandler.GetUserRoles)
			r.Get("/roles", authHandler.GetRoles)
		})
	})
}

func registerCategoryRoutes(r chi.Router, authService auth.AuthService, categoryHandler *category.CategoryHandler) {
	r.Route("/categories", func(r chi.Router) {

		r.Get("/", categoryHandler.GetAll)
		r.Get("/{id}", categoryHandler.GetById)
		r.Get("/code/{code}", categoryHandler.GetByCode)
		r.Get("/category/{category_id}", categoryHandler.GetByCategoryID)
		r.Post("/search", categoryHandler.GetWithSearchCriteria)

		// Защищенные маршруты (требуют аутентификации)
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(authService, "admin", "guest"))

			r.Post("/", categoryHandler.Create)
			r.Delete("/{id}", categoryHandler.Delete)
		})
	})
}

func registerProductRoutes(r chi.Router, authService auth.AuthService, productHandler *product.ProductHandler) {
	r.Route("/products", func(r chi.Router) {

		r.Get("/", productHandler.GetAll)
		r.Get("/{id}", productHandler.GetById)
		r.Get("/code/{code}", productHandler.GetByCode)
		r.Get("/category/{category_id}", productHandler.GetByCategoryID)
		r.Post("/search", productHandler.GetWithSearchCriteria)

		// Защищенные маршруты (требуют аутентификации)
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(authService, "admin", "guest"))

			r.Post("/", productHandler.Create)
			r.Delete("/{id}", productHandler.Delete)
		})

	})
}

// setupStaticRoutes настраивает раздачу статических файлов
func setupStaticRoutes(r chi.Router, staticFilesPath string) {
	// Создаем файловый сервер для статических файлов
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, staticFilesPath))

	// Настраиваем маршрут для статических файлов
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(filesDir)))
}

// registerProductMediaRoutes регистрирует маршруты для работы с медиа товаров
func registerProductMediaRoutes(r chi.Router, authService auth.AuthService, productMediaHandler *productmedia.ProductMediaHandler) {
	r.Route("/product-media", func(r chi.Router) {
		r.Get("/product/{product_id}", productMediaHandler.GetByProductID)
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(authService, "admin", "guest"))

			r.Post("/upload", productMediaHandler.UploadImage)
			r.Delete("/{id}", productMediaHandler.Delete)
		})
	})
}

func registerCartRoutes(r chi.Router, authService auth.AuthService, cartHandler *cart.CartHandler) {
	r.Route("/carts", func(r chi.Router) {
		r.Use(AuthMiddleware(authService, "admin", "guest"))

		r.Get("/{id}", cartHandler.GetById)
		r.Get("/person/{person_id}", cartHandler.GetByPersonID)
		r.Post("/search", cartHandler.GetWithSearchCriteria)

		r.Post("/", cartHandler.Create)
	})
}

func registerCartItemRoutes(r chi.Router, authService auth.AuthService, cartItemHandler *cartitem.CartItemHandler) {
	r.Route("/cart-items", func(r chi.Router) {
		r.Use(AuthMiddleware(authService, "admin", "guest"))

		r.Get("/{id}", cartItemHandler.GetById)
		r.Get("/cart/{cart_id}", cartItemHandler.GetByCartID)
		r.Post("/search", cartItemHandler.GetWithSearchCriteria)

		r.Post("/", cartItemHandler.Create)
		r.Patch("/", cartItemHandler.Update)
		r.Delete("/{id}", cartItemHandler.Delete)
	})
}

func registerOrderRoutes(r chi.Router, authService auth.AuthService, orderHandler *order.OrderHandler) {
	r.Route("/orders", func(r chi.Router) {
		r.Use(AuthMiddleware(authService, "admin", "guest"))

		r.Get("/{id}", orderHandler.GetById)
		r.Get("/client/{client_id}", orderHandler.GetByClientID)
		r.Get("/manager/{manager_id}", orderHandler.GetByManagerID)
		r.Get("/manager/{manager_id}/status/{status}", orderHandler.GetByManagerIDAndStatus)
		r.Get("/status/{status}", orderHandler.GetByStatus)
		r.Post("/search", orderHandler.GetWithSearchCriteria)

		r.Post("/", orderHandler.Create)
		r.Delete("/{id}", orderHandler.Delete)

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(authService, "admin"))

			r.Post("/{id}/change-status/{status}", orderHandler.ChangeStatus)
		})
	})
}

func registerOrderItemRoutes(r chi.Router, authService auth.AuthService, orderItemHandler *orderitem.OrderItemHandler) {
	r.Route("/order-items", func(r chi.Router) {
		r.Use(AuthMiddleware(authService, "admin", "guest"))

		r.Get("/{id}", orderItemHandler.GetById)
		r.Get("/order/{order_id}", orderItemHandler.GetByOrderID)
		r.Post("/search", orderItemHandler.GetWithSearchCriteria)

		r.Post("/", orderItemHandler.Create)
		r.Delete("/{id}", orderItemHandler.Delete)

		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(authService, "admin"))

			r.Post("/{id}/change-status/{status}", orderItemHandler.ChangeStatus)
		})
	})
}

func RegisterSwaggerRoutes(router chi.Router) {
	// Настройка Swagger
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // URL для JSON документации
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)

	router.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	router.Get("/swagger/*", swaggerHandler)
}
