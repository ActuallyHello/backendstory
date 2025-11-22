package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/ActuallyHello/backendstory/internal/core/container"
	handlers "github.com/ActuallyHello/backendstory/internal/server/handlers"
	appMiddleware "github.com/ActuallyHello/backendstory/internal/server/middleware"
	"github.com/ActuallyHello/backendstory/internal/services/auth"
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
	r.Use(appMiddleware.ErrorHandler)
	// static files
	setupStaticRoutes(r, staticFilesPath)

	r.Route(apiV1, func(r chi.Router) {
		r.Post(register, container.GetAuthHandler().Register)
		r.Post(login, container.GetAuthHandler().Login)

		registerAuthRoutes(r, container.GetAuthService(), container.GetAuthHandler())

		registerEnumRoutes(r, container.GetAuthService(), container.GetEnumHandler())
		registerEnumValuesRoutes(r, container.GetAuthService(), container.GetEnumValueHandler())
		registerPersonRoutes(r, container.GetAuthService(), container.GetPersonHandler())
		registerCategoryRoutes(r, container.GetAuthService(), container.GetCategoryHandler())
		registerProductRoutes(r, container.GetAuthService(), container.GetProductHandler())
		registerProductMediaRoutes(r, container.GetAuthService(), container.GetProductMediaHandler())
		registerCartRoutes(r, container.GetAuthService(), container.GetCartHandler())
		registerCartItemRoutes(r, container.GetAuthService(), container.GetCartItemHandler())
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

func registerEnumRoutes(r chi.Router, authService auth.AuthService, enumHandler *handlers.EnumHandler) {
	// r.Use(appMiddleware.KeycloakAuthMiddleware(kc, "admin"))
	r.Route("/enumerations", func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(authService, "admin"))

		r.Get("/", enumHandler.GetAll)
		r.Get(byId, enumHandler.GetById)
		r.Get("/code/{code}", enumHandler.GetByCode)
		r.Post("/search", enumHandler.GetWithSearchCriteria)

		r.Post("/", enumHandler.Create)
		r.Delete(byId, enumHandler.Delete)
	})
}

func registerEnumValuesRoutes(r chi.Router, authService auth.AuthService, enumValueHandler *handlers.EnumValueHandler) {
	r.Route("/enumeration-values", func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(authService, "admin"))

		r.Get("/", enumValueHandler.GetAll)
		r.Get(byId, enumValueHandler.GetById)
		r.Get("/enumeration/{enumeration_id}", enumValueHandler.GetByEnumId)
		r.Post("/search", enumValueHandler.GetWithSearchCriteria)

		r.Post("/", enumValueHandler.Create)
		r.Delete(byId, enumValueHandler.Delete)
	})
}

func registerPersonRoutes(r chi.Router, authService auth.AuthService, personHandler *handlers.PersonHandler) {
	r.Route("/persons", func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(authService, "admin"))

		r.Get("/", personHandler.GetAll)
		r.Get(byId, personHandler.GetById)
		r.Get("/user/{user_login}", personHandler.GetByUserLogin)
		r.Post("/search", personHandler.GetWithSearchCriteria)

		r.Post("/", personHandler.Create)
		r.Delete(byId, personHandler.Delete)
	})
}

func registerAuthRoutes(r chi.Router, authService auth.AuthService, authHandler *handlers.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(authService, "admin", "guest"))

		r.Get("/users", authHandler.GetUsers)
		r.Get("/users/{username}", authHandler.GetUser)
		r.Get("/users/{username}/roles", authHandler.GetUserRoles)
		r.Get("/roles", authHandler.GetRoles)

		r.Get("/token", authHandler.GetHeaderTokenInfo)
		r.Post("/token", authHandler.GetBodyTokenInfo)
	})
}

func registerCategoryRoutes(r chi.Router, authService auth.AuthService, categoryHandler *handlers.CategoryHandler) {
	r.Route("/categories", func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(authService, "admin", "guest"))

		r.Get("/", categoryHandler.GetAll)
		r.Get("/{id}", categoryHandler.GetById)
		r.Get("/code/{code}", categoryHandler.GetByCode)
		r.Get("/category/{category_id}", categoryHandler.GetByCategoryID)
		r.Post("/search", categoryHandler.GetWithSearchCriteria)

		r.Post("/", categoryHandler.Create)
		r.Delete("/{id}", categoryHandler.Delete)
	})
}

func registerProductRoutes(r chi.Router, authService auth.AuthService, productHandler *handlers.ProductHandler) {
	r.Route("/products", func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(authService, "admin", "guest"))

		r.Get("/", productHandler.GetAll)
		r.Get("/{id}", productHandler.GetById)
		r.Get("/code/{code}", productHandler.GetByCode)
		r.Get("/category/{category_id}", productHandler.GetByCategoryID)
		r.Post("/search", productHandler.GetWithSearchCriteria)

		r.Post("/", productHandler.Create)
		r.Delete("/{id}", productHandler.Delete)
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
func registerProductMediaRoutes(r chi.Router, authService auth.AuthService, productMediaHandler *handlers.ProductMediaHandler) {
	r.Route("/product-media", func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(authService, "admin", "guest"))

		// Загрузка изображения
		r.Post("/upload", productMediaHandler.UploadImage)

		// Получение изображений по product_id
		r.Get("/product/{product_id}", productMediaHandler.GetByProductID)

		// Удаление изображения
		r.Delete("/{id}", productMediaHandler.Delete)
	})
}

func registerCartRoutes(r chi.Router, authService auth.AuthService, cartHandler *handlers.CartHandler) {
	r.Route("/carts", func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(authService, "admin", "guest"))

		r.Get("/{id}", cartHandler.GetById)
		r.Get("/person/{person_id}", cartHandler.GetByPersonID)
		r.Post("/search", cartHandler.GetWithSearchCriteria)

		r.Post("/", cartHandler.Create)
	})
}

func registerCartItemRoutes(r chi.Router, authService auth.AuthService, cartItemHandler *handlers.CartItemHandler) {
	r.Route("/cart-items", func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(authService, "admin", "guest"))

		r.Get("/{id}", cartItemHandler.GetById)
		r.Get("/cart/{cart_id}", cartItemHandler.GetByCartID)
		r.Post("/search", cartItemHandler.GetWithSearchCriteria)

		r.Post("/", cartItemHandler.Create)
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
