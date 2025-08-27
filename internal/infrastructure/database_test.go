package infrastructure

import (
	"testing"
)

func TestNewDatabase(t *testing.T) {
	tests := []struct {
		name      string
		config    DatabaseConfig
		expectErr bool
	}{
		{
			name: "Valid SQLite config",
			config: DatabaseConfig{
				Driver: "sqlite3",
				DSN:    ":memory:",
			},
			expectErr: false,
		},
		{
			name: "Invalid driver",
			config: DatabaseConfig{
				Driver: "invalid_driver",
				DSN:    ":memory:",
			},
			expectErr: true,
		},
		{
			name: "Invalid DSN for postgres",
			config: DatabaseConfig{
				Driver: "postgres",
				DSN:    "invalid_dsn",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDatabase(tt.config)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if db == nil {
				t.Error("Expected database connection but got nil")
				return
			}

			// Test that we can ping the database
			err = db.Ping()
			if err != nil {
				t.Errorf("Failed to ping database: %v", err)
			}

			// Clean up
			db.Close()
		})
	}
}

func TestDatabaseConfig(t *testing.T) {
	config := DatabaseConfig{
		Driver: "sqlite3",
		DSN:    "test.db",
	}

	if config.Driver != "sqlite3" {
		t.Errorf("Expected driver sqlite3, got %s", config.Driver)
	}

	if config.DSN != "test.db" {
		t.Errorf("Expected DSN test.db, got %s", config.DSN)
	}
}

func TestDatabase_Initialization(t *testing.T) {
	config := DatabaseConfig{
		Driver: "sqlite3",
		DSN:    ":memory:",
	}

	db, err := NewDatabase(config)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test that the users table was created
	query := "SELECT name FROM sqlite_master WHERE type='table' AND name='users'"
	var tableName string
	err = db.QueryRow(query).Scan(&tableName)
	if err != nil {
		t.Errorf("Users table was not created: %v", err)
	}

	if tableName != "users" {
		t.Errorf("Expected table name 'users', got '%s'", tableName)
	}
}

func TestDatabase_TableStructure(t *testing.T) {
	config := DatabaseConfig{
		Driver: "sqlite3",
		DSN:    ":memory:",
	}

	db, err := NewDatabase(config)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test table structure by checking columns
	query := "PRAGMA table_info(users)"
	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("Failed to get table info: %v", err)
	}
	defer rows.Close()

	expectedColumns := map[string]bool{
		"id":         false,
		"email":      false,
		"password":   false,
		"firstname":  false,
		"lastname":   false,
		"phone":      false,
		"birthday":   false,
		"created_at": false,
		"updated_at": false,
	}

	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, primaryKey int
		var defaultValue interface{}

		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey)
		if err != nil {
			t.Errorf("Failed to scan row: %v", err)
			continue
		}

		if _, exists := expectedColumns[name]; exists {
			expectedColumns[name] = true
		}
	}

	// Check that all expected columns were found
	for column, found := range expectedColumns {
		if !found {
			t.Errorf("Expected column '%s' not found in users table", column)
		}
	}
}
