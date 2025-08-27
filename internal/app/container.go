package app

import (
	"database/sql"

	"hello-world/internal/domain"
	"hello-world/internal/infrastructure"
	"hello-world/internal/interfaces"
	"hello-world/internal/usecase"
	"hello-world/pkg/config"
)

// Container holds all the application dependencies
type Container struct {
	Config      *config.Config
	Database    *sql.DB
	UserRepo    domain.UserRepository
	AuthService domain.AuthService
	UserService domain.UserService
	Router      *interfaces.Router
}

// NewContainer creates and wires all dependencies
func NewContainer() (*Container, error) {
	// Load configuration
	cfg := config.Load()

	// Initialize infrastructure layer
	db, err := infrastructure.NewDatabase(infrastructure.DatabaseConfig{
		Driver: cfg.Database.Driver,
		DSN:    cfg.Database.DSN,
	})
	if err != nil {
		return nil, err
	}

	// Initialize repositories (adapters)
	userRepo := infrastructure.NewSQLiteUserRepository(db)

	// Initialize services (adapters)
	authService := infrastructure.NewJWTAuthService()

	// Initialize use cases (application layer)
	userService := usecase.NewUserUseCase(userRepo, authService)

	// Initialize interface layer
	router := interfaces.NewRouter(userService, authService)

	return &Container{
		Config:      cfg,
		Database:    db,
		UserRepo:    userRepo,
		AuthService: authService,
		UserService: userService,
		Router:      router,
	}, nil
}

// Close cleans up resources
func (c *Container) Close() error {
	if c.Database != nil {
		return c.Database.Close()
	}
	return nil
}
