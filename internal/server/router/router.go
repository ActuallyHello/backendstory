package router

import (
	"net/http"

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

func SetupRouter(container *container.AppContainer) http.Handler {
	r := chi.NewRouter()

	// TODO:
	//  3. unless-stop -> 3 retry

	r.Use(middleware.Logger)
	// custom ErrorHandler
	r.Use(appMiddleware.ErrorHandler)

	r.Route(apiV1, func(r chi.Router) {
		r.Post(register, container.GetAuthHandler().Register)
		r.Post(login, container.GetAuthHandler().Login)

		registerEnumRoutes(r, container.GetAuthService(), container.GetEnumHandler())
		registerEnumValuesRoutes(r, container.GetAuthService(), container.GetEnumValueHandler())

		registerAuthRoutes(r, container.GetAuthService(), container.GetAuthHandler())
		registerPersonRoutes(r, container.GetAuthService(), container.GetPersonHandler())
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
		// r.Use(appMiddleware.AuthMiddleware(authService, "admin"))

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

		r.Post("/", personHandler.Create)
		r.Delete(byId, personHandler.Delete)
	})
}

func registerAuthRoutes(r chi.Router, authService auth.AuthService, authHandler *handlers.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(authService, "admin", "guest"))

		r.Get("/users", authHandler.GetUsers)
		r.Get("/roles", authHandler.GetRoles)
		r.Get("/users/{username}/roles", authHandler.GetUserRoles)
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
