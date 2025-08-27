package app

import (
	"testing"

	"hello-world/internal/domain"
)

func TestNewContainer(t *testing.T) {
	container, err := NewContainer()
	if err != nil {
		t.Errorf("Unexpected error creating container: %v", err)
	}

	if container == nil {
		t.Error("Expected container to be created but got nil")
	}

	// Test that all dependencies are properly initialized
	if container.Config == nil {
		t.Error("Expected config to be initialized")
	}

	if container.Database == nil {
		t.Error("Expected database to be initialized")
	}

	if container.UserRepo == nil {
		t.Error("Expected user repository to be initialized")
	}

	if container.AuthService == nil {
		t.Error("Expected auth service to be initialized")
	}

	if container.UserService == nil {
		t.Error("Expected user service to be initialized")
	}

	if container.Router == nil {
		t.Error("Expected router to be initialized")
	}

	// Test that dependencies implement the correct interfaces
	var _ domain.UserRepository = container.UserRepo
	var _ domain.AuthService = container.AuthService
	var _ domain.UserService = container.UserService

	// Clean up
	container.Close()
}

func TestContainer_Close(t *testing.T) {
	container, err := NewContainer()
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	err = container.Close()
	if err != nil {
		t.Errorf("Unexpected error closing container: %v", err)
	}
}

func TestContainer_DatabaseConnection(t *testing.T) {
	container, err := NewContainer()
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}
	defer container.Close()

	// Test that database connection is working
	err = container.Database.Ping()
	if err != nil {
		t.Errorf("Database connection failed: %v", err)
	}
}

func TestContainer_ConfigValidation(t *testing.T) {
	container, err := NewContainer()
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}
	defer container.Close()

	// Test basic config validation
	if container.Config.Server.Host == "" {
		t.Error("Expected server host to be configured")
	}

	if container.Config.Server.Port == "" {
		t.Error("Expected server port to be configured")
	}

	if container.Config.Database.Driver == "" {
		t.Error("Expected database driver to be configured")
	}

	if container.Config.Database.DSN == "" {
		t.Error("Expected database DSN to be configured")
	}
}
