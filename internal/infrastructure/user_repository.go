package infrastructure

import (
	"context"
	"database/sql"
	"strings"

	"hello-world/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteUserRepository implements domain.UserRepository using SQLite
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new SQLite user repository
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// Create inserts a new user into the database
func (r *SQLiteUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password, firstname, lastname, phone, birthday, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email, user.Password, user.FirstName, user.LastName,
		user.Phone, user.Birthday, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return domain.ErrUserAlreadyExists
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

// GetByEmail retrieves a user by email
func (r *SQLiteUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password, firstname, lastname, phone, birthday, created_at, updated_at 
		FROM users WHERE email = ?
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Phone, &user.Birthday, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *SQLiteUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT id, email, password, firstname, lastname, phone, birthday, created_at, updated_at 
		FROM users WHERE id = ?
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Phone, &user.Birthday, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// Update updates an existing user
func (r *SQLiteUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET 
			email = ?, password = ?, firstname = ?, lastname = ?, 
			phone = ?, birthday = ?, updated_at = ? 
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.Password, user.FirstName, user.LastName,
		user.Phone, user.Birthday, user.UpdatedAt, user.ID,
	)

	return err
}

// Delete removes a user by ID
func (r *SQLiteUserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Exists checks if a user with the given email exists
func (r *SQLiteUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
