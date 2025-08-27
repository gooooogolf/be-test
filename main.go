package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "hello-world/docs" // Import generated swagger docs
)

// @title Go Chi API
// @version 1.0
// @description A simple Go API with user authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3333
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

type Response struct {
	Message string `json:"message"`
}

func main() {
	// Initialize database
	InitDB()
	defer CloseDB()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Public routes
	r.Get("/", HelloHandler)
	r.Post("/register", RegisterHandler)
	r.Post("/login", LoginHandler)

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3333/swagger/doc.json"),
	))

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Get("/me", MeHandler)
	})

	log.Println("Server starting on :3333")
	log.Println("Swagger UI available at: http://localhost:3333/swagger/")
	http.ListenAndServe(":3333", r)
}
