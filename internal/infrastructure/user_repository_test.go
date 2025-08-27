package infrastructure

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"hello-world/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create users table
	createTable := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		firstname TEXT NOT NULL,
		lastname TEXT NOT NULL,
		phone TEXT NOT NULL,
		birthday DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`

	_, err = db.Exec(createTable)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	return db
}

func TestNewSQLiteUserRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	if repo == nil {
		t.Error("Expected repository to be created but got nil")
	}
}

func TestSQLiteUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	user := &domain.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set after creation")
	}
}

func TestSQLiteUserRepository_Create_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	user1 := &domain.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user2 := &domain.User{
		Email:     "test@example.com", // Same email
		Password:  "hashedpassword2",
		FirstName: "Jane",
		LastName:  "Smith",
		Phone:     "0987654321",
		Birthday:  time.Date(1991, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create first user
	err := repo.Create(ctx, user1)
	if err != nil {
		t.Errorf("Unexpected error creating first user: %v", err)
	}

	// Try to create second user with same email
	err = repo.Create(ctx, user2)
	if err != domain.ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestSQLiteUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get user by email
	foundUser, err := repo.GetByEmail(ctx, "test@example.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if foundUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, foundUser.Email)
	}
	if foundUser.FirstName != user.FirstName {
		t.Errorf("Expected first name %s, got %s", user.FirstName, foundUser.FirstName)
	}
}

func TestSQLiteUserRepository_GetByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	_, err := repo.GetByEmail(ctx, "nonexistent@example.com")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

func TestSQLiteUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Get user by ID
	foundUser, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if foundUser.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, foundUser.ID)
	}
	if foundUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, foundUser.Email)
	}
}

func TestSQLiteUserRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 999)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

func TestSQLiteUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Update user
	user.FirstName = "Jane"
	user.LastName = "Smith"
	user.UpdatedAt = time.Now()

	err = repo.Update(ctx, user)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify update
	foundUser, err := repo.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if foundUser.FirstName != "Jane" {
		t.Errorf("Expected first name Jane, got %s", foundUser.FirstName)
	}
	if foundUser.LastName != "Smith" {
		t.Errorf("Expected last name Smith, got %s", foundUser.LastName)
	}
}

func TestSQLiteUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Delete user
	err = repo.Delete(ctx, user.ID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, user.ID)
	if err == nil {
		t.Error("Expected error for deleted user")
	}
}

func TestSQLiteUserRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLiteUserRepository(db)
	ctx := context.Background()

	// Check non-existent user
	exists, err := repo.Exists(ctx, "nonexistent@example.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if exists {
		t.Error("Expected user to not exist")
	}

	// Create a user
	user := &domain.User{
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Check existing user
	exists, err = repo.Exists(ctx, "test@example.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !exists {
		t.Error("Expected user to exist")
	}
}
