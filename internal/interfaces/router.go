package interfaces

import (
	"hello-world/internal/domain"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router holds all the route handlers and dependencies
type Router struct {
	userHandler    *UserHandler
	authMiddleware *AuthMiddleware
}

// NewRouter creates a new router with all dependencies
func NewRouter(
	userService domain.UserService,
	authService domain.AuthService,
) *Router {
	return &Router{
		userHandler:    NewUserHandler(userService),
		authMiddleware: NewAuthMiddleware(authService),
	}
}

// SetupRoutes configures all the routes
func (router *Router) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Public routes
	r.Get("/", router.userHandler.HelloHandler)
	r.Post("/register", router.userHandler.RegisterHandler)
	r.Post("/login", router.userHandler.LoginHandler)

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3333/swagger/doc.json"),
	))

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(router.authMiddleware.Middleware)
		r.Get("/me", router.userHandler.MeHandler)
	})

	return r
}
