package main

import (
	"log"
	"net/http"

	"hello-world/internal/app"

	_ "hello-world/docs" // Import generated swagger docs
)

// @title Go Chi API
// @version 1.0
// @description A simple Go API with user authentication using Clean Architecture
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

func main() {
	// Initialize dependency injection container
	container, err := app.NewContainer()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}
	defer container.Close()

	// Setup routes
	r := container.Router.SetupRoutes()

	// Start server
	serverAddr := container.Config.Server.Host + ":" + container.Config.Server.Port
	log.Printf("Server starting on %s", serverAddr)
	log.Printf("Swagger UI available at: http://localhost:%s/swagger/", container.Config.Server.Port)

	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
