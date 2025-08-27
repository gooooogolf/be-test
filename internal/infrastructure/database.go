package infrastructure

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver string
	DSN    string
}

// NewDatabase creates a new database connection
func NewDatabase(config DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open(config.Driver, config.DSN)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Create tables
	if err = createTables(db); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// createTables creates the necessary database tables
func createTables(db *sql.DB) error {
	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			firstname TEXT NOT NULL,
			lastname TEXT NOT NULL,
			phone TEXT NOT NULL,
			birthday DATE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(createUsersTable)
	if err != nil {
		return err
	}

	// Create indexes for better performance
	createIndexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
		`CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);`,
	}

	for _, indexSQL := range createIndexes {
		_, err = db.Exec(indexSQL)
		if err != nil {
			return err
		}
	}

	return nil
}
