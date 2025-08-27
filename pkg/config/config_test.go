package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test default configuration
	config := Load()

	if config == nil {
		t.Error("Expected config to be loaded but got nil")
	}

	// Test default values
	if config.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default host 0.0.0.0, got %s", config.Server.Host)
	}

	if config.Server.Port != "3333" {
		t.Errorf("Expected default port 3333, got %s", config.Server.Port)
	}

	if config.Database.Driver != "sqlite3" {
		t.Errorf("Expected default driver sqlite3, got %s", config.Database.Driver)
	}

	if config.Database.DSN != "./app.db" {
		t.Errorf("Expected default DSN ./app.db, got %s", config.Database.DSN)
	}
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	originalHost := os.Getenv("SERVER_HOST")
	originalPort := os.Getenv("SERVER_PORT")
	originalDriver := os.Getenv("DB_DRIVER")
	originalDSN := os.Getenv("DB_DSN")

	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("DB_DSN", "postgres://user:pass@localhost/db")

	// Clean up after test
	defer func() {
		os.Setenv("SERVER_HOST", originalHost)
		os.Setenv("SERVER_PORT", originalPort)
		os.Setenv("DB_DRIVER", originalDriver)
		os.Setenv("DB_DSN", originalDSN)
	}()

	config := Load()

	if config.Server.Host != "127.0.0.1" {
		t.Errorf("Expected host from env 127.0.0.1, got %s", config.Server.Host)
	}

	if config.Server.Port != "8080" {
		t.Errorf("Expected port from env 8080, got %s", config.Server.Port)
	}

	if config.Database.Driver != "postgres" {
		t.Errorf("Expected driver from env postgres, got %s", config.Database.Driver)
	}

	if config.Database.DSN != "postgres://user:pass@localhost/db" {
		t.Errorf("Expected DSN from env, got %s", config.Database.DSN)
	}
}

func TestConfig_Structure(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: "9090",
		},
		Database: DatabaseConfig{
			Driver: "mysql",
			DSN:    "mysql://localhost/test",
		},
	}

	if config.Server.Host != "localhost" {
		t.Errorf("Expected host localhost, got %s", config.Server.Host)
	}

	if config.Server.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", config.Server.Port)
	}

	if config.Database.Driver != "mysql" {
		t.Errorf("Expected driver mysql, got %s", config.Database.Driver)
	}

	if config.Database.DSN != "mysql://localhost/test" {
		t.Errorf("Expected DSN mysql://localhost/test, got %s", config.Database.DSN)
	}
}

func TestServerConfig(t *testing.T) {
	serverConfig := ServerConfig{
		Host: "192.168.1.1",
		Port: "4444",
	}

	if serverConfig.Host != "192.168.1.1" {
		t.Errorf("Expected host 192.168.1.1, got %s", serverConfig.Host)
	}

	if serverConfig.Port != "4444" {
		t.Errorf("Expected port 4444, got %s", serverConfig.Port)
	}
}

func TestDatabaseConfig(t *testing.T) {
	dbConfig := DatabaseConfig{
		Driver: "sqlite3",
		DSN:    "test.db",
	}

	if dbConfig.Driver != "sqlite3" {
		t.Errorf("Expected driver sqlite3, got %s", dbConfig.Driver)
	}

	if dbConfig.DSN != "test.db" {
		t.Errorf("Expected DSN test.db, got %s", dbConfig.DSN)
	}
}

func TestLoad_EmptyEnvironmentVariables(t *testing.T) {
	// Set empty environment variables
	originalHost := os.Getenv("SERVER_HOST")
	originalPort := os.Getenv("SERVER_PORT")
	originalDriver := os.Getenv("DB_DRIVER")
	originalDSN := os.Getenv("DB_DSN")

	os.Setenv("SERVER_HOST", "")
	os.Setenv("SERVER_PORT", "")
	os.Setenv("DB_DRIVER", "")
	os.Setenv("DB_DSN", "")

	// Clean up after test
	defer func() {
		os.Setenv("SERVER_HOST", originalHost)
		os.Setenv("SERVER_PORT", originalPort)
		os.Setenv("DB_DRIVER", originalDriver)
		os.Setenv("DB_DSN", originalDSN)
	}()

	config := Load()

	// Should fall back to defaults
	if config.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default host 0.0.0.0, got %s", config.Server.Host)
	}

	if config.Server.Port != "3333" {
		t.Errorf("Expected default port 3333, got %s", config.Server.Port)
	}

	if config.Database.Driver != "sqlite3" {
		t.Errorf("Expected default driver sqlite3, got %s", config.Database.Driver)
	}

	if config.Database.DSN != "./app.db" {
		t.Errorf("Expected default DSN ./app.db, got %s", config.Database.DSN)
	}
}
