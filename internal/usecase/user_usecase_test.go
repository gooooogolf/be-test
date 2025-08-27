package usecase

import (
	"context"
	"testing"
	"time"

	"hello-world/internal/domain"
)

// MockUserRepository implements domain.UserRepository for testing
type MockUserRepository struct {
	users  map[string]*domain.User
	nextID int
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[string]*domain.User),
		nextID: 1,
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = m.nextID
	m.nextID++
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return domain.ErrUserNotFound
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	_, exists := m.users[email]
	return exists, nil
}

// MockAuthService implements domain.AuthService for testing
type MockAuthService struct{}

func NewMockAuthService() *MockAuthService {
	return &MockAuthService{}
}

func (m *MockAuthService) GenerateToken(userID int, email string) (string, error) {
	return "mock_token", nil
}

func (m *MockAuthService) ValidateToken(token string) (*domain.TokenClaims, error) {
	return &domain.TokenClaims{UserID: 1, Email: "test@example.com"}, nil
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	return "hashed_" + password, nil
}

func (m *MockAuthService) ComparePassword(hashedPassword, password string) error {
	if hashedPassword == "hashed_"+password {
		return nil
	}
	return domain.ErrInvalidCredentials
}

// Test functions
func TestUserUseCase_Register(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := NewMockAuthService()
	userService := NewUserUseCase(userRepo, authService)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	firstName := "John"
	lastName := "Doe"
	phone := "1234567890"
	birthday := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	// Act
	user, err := userService.Register(ctx, email, password, firstName, lastName, phone, birthday)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.Email != email {
		t.Fatalf("Expected email %s, got %s", email, user.Email)
	}
	if user.FirstName != firstName {
		t.Fatalf("Expected first name %s, got %s", firstName, user.FirstName)
	}
	if user.ID == 0 {
		t.Fatal("Expected user ID to be set")
	}
	if user.Password != "" {
		t.Fatal("Expected password to be cleared from response")
	}
}

func TestUserUseCase_Login(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := NewMockAuthService()
	userService := NewUserUseCase(userRepo, authService)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	// Create a user first
	_, err := userService.Register(ctx, email, password, "John", "Doe", "1234567890", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Check that user was created with hashed password in the repository
	storedUser, err := userRepo.GetByEmail(ctx, email)
	if err != nil {
		t.Fatalf("Failed to get stored user: %v", err)
	}
	t.Logf("Stored user password: %s", storedUser.Password)

	// Act
	token, loginUser, err := userService.Login(ctx, email, password)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("Expected token to be returned")
	}
	if loginUser.Email != email {
		t.Fatalf("Expected email %s, got %s", email, loginUser.Email)
	}
	if loginUser.Password != "" {
		t.Fatal("Expected password to be cleared from response")
	}
}

func TestUserUseCase_Register_DuplicateEmail(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := NewMockAuthService()
	userService := NewUserUseCase(userRepo, authService)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	firstName := "John"
	lastName := "Doe"
	phone := "1234567890"
	birthday := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	// Register user first time
	_, err := userService.Register(ctx, email, password, firstName, lastName, phone, birthday)
	if err != nil {
		t.Fatalf("Failed to register user first time: %v", err)
	}

	// Act - try to register same user again
	_, err = userService.Register(ctx, email, password, firstName, lastName, phone, birthday)

	// Assert
	if err == nil {
		t.Fatal("Expected error for duplicate email")
	}
	if err != domain.ErrUserAlreadyExists {
		t.Fatalf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := NewMockAuthService()
	userService := NewUserUseCase(userRepo, authService)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	// Create a user first
	_, err := userService.Register(ctx, email, password, "John", "Doe", "1234567890", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Get the created user
	user, err := userRepo.GetByEmail(ctx, email)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Act - Update user
	newBirthday := time.Date(1995, 5, 15, 0, 0, 0, 0, time.UTC)
	updatedUser, err := userService.UpdateUser(ctx, user.ID, "Jane", "Smith", "0987654321", &newBirthday)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if updatedUser.FirstName != "Jane" {
		t.Fatalf("Expected first name Jane, got %s", updatedUser.FirstName)
	}
	if updatedUser.LastName != "Smith" {
		t.Fatalf("Expected last name Smith, got %s", updatedUser.LastName)
	}
	if updatedUser.Phone != "0987654321" {
		t.Fatalf("Expected phone 0987654321, got %s", updatedUser.Phone)
	}
	if !updatedUser.Birthday.Equal(newBirthday) {
		t.Fatalf("Expected birthday %v, got %v", newBirthday, updatedUser.Birthday)
	}
	if updatedUser.Password != "" {
		t.Fatal("Expected password to be cleared from response")
	}
}

func TestUserUseCase_UpdateUser_NotFound(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := NewMockAuthService()
	userService := NewUserUseCase(userRepo, authService)

	ctx := context.Background()

	// Act - Try to update non-existent user
	_, err := userService.UpdateUser(ctx, 999, "Jane", "Smith", "0987654321", nil)

	// Assert
	if err != domain.ErrUserNotFound {
		t.Fatalf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUserUseCase_UpdateUser_PartialUpdate(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := NewMockAuthService()
	userService := NewUserUseCase(userRepo, authService)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	// Create a user first
	originalUser, err := userService.Register(ctx, email, password, "John", "Doe", "1234567890", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Act - Update only first name
	updatedUser, err := userService.UpdateUser(ctx, originalUser.ID, "Jane", "", "", nil)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if updatedUser.FirstName != "Jane" {
		t.Fatalf("Expected first name Jane, got %s", updatedUser.FirstName)
	}
	// Other fields should remain unchanged
	if updatedUser.LastName != "Doe" {
		t.Fatalf("Expected last name to remain Doe, got %s", updatedUser.LastName)
	}
	if updatedUser.Phone != "1234567890" {
		t.Fatalf("Expected phone to remain 1234567890, got %s", updatedUser.Phone)
	}
}

func TestUserUseCase_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := NewMockAuthService()
	userService := NewUserUseCase(userRepo, authService)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	// Create a user first
	_, err := userService.Register(ctx, email, password, "John", "Doe", "1234567890", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Act - Try to login with wrong password
	_, _, err = userService.Login(ctx, email, "wrongpassword")

	// Assert
	if err != domain.ErrInvalidCredentials {
		t.Fatalf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUserUseCase_Login_UserNotFound(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := NewMockAuthService()
	userService := NewUserUseCase(userRepo, authService)

	ctx := context.Background()

	// Act - Try to login with non-existent user
	_, _, err := userService.Login(ctx, "nonexistent@example.com", "password")

	// Assert
	if err != domain.ErrInvalidCredentials {
		t.Fatalf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUserUseCase_GetUserProfile(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := NewMockAuthService()
	userService := NewUserUseCase(userRepo, authService)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	// Create a user first
	createdUser, err := userService.Register(ctx, email, password, "John", "Doe", "1234567890", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Act
	user, err := userService.GetUserProfile(ctx, createdUser.ID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.Email != email {
		t.Fatalf("Expected email %s, got %s", email, user.Email)
	}
	if user.Password != "" {
		t.Fatal("Expected password to be cleared from response")
	}
}

func TestUserUseCase_Register_HashPasswordError(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := &MockAuthServiceWithError{hashError: true}
	userService := NewUserUseCase(userRepo, authService)

	ctx := context.Background()

	// Act
	_, err := userService.Register(ctx, "test@example.com", "password", "John", "Doe", "1234567890", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))

	// Assert
	if err != domain.ErrPasswordHashError {
		t.Fatalf("Expected ErrPasswordHashError, got %v", err)
	}
}

func TestUserUseCase_Register_CreateUserError(t *testing.T) {
	// Arrange
	userRepo := &MockUserRepositoryWithError{createError: true}
	authService := NewMockAuthService()
	userService := NewUserUseCase(userRepo, authService)

	ctx := context.Background()

	// Act
	_, err := userService.Register(ctx, "test@example.com", "password", "John", "Doe", "1234567890", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))

	// Assert
	if err != domain.ErrUserCreationError {
		t.Fatalf("Expected ErrUserCreationError, got %v", err)
	}
}

// Mock implementations with error scenarios
type MockAuthServiceWithError struct {
	hashError     bool
	tokenError    bool
	compareError  bool
	validateError bool
}

func (m *MockAuthServiceWithError) GenerateToken(userID int, email string) (string, error) {
	if m.tokenError {
		return "", domain.ErrTokenGenerationError
	}
	return "mock_token", nil
}

func (m *MockAuthServiceWithError) ValidateToken(token string) (*domain.TokenClaims, error) {
	if m.validateError {
		return nil, domain.ErrInvalidToken
	}
	return &domain.TokenClaims{UserID: 1, Email: "test@example.com"}, nil
}

func (m *MockAuthServiceWithError) HashPassword(password string) (string, error) {
	if m.hashError {
		return "", domain.ErrPasswordHashError
	}
	return "hashed_" + password, nil
}

func (m *MockAuthServiceWithError) ComparePassword(hashedPassword, password string) error {
	if m.compareError {
		return domain.ErrInvalidCredentials
	}
	if hashedPassword == "hashed_"+password {
		return nil
	}
	return domain.ErrInvalidCredentials
}

type MockUserRepositoryWithError struct {
	createError bool
	existsError bool
}

func (m *MockUserRepositoryWithError) Create(ctx context.Context, user *domain.User) error {
	if m.createError {
		return domain.ErrUserCreationError
	}
	return nil
}

func (m *MockUserRepositoryWithError) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepositoryWithError) GetByID(ctx context.Context, id int) (*domain.User, error) {
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepositoryWithError) Update(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *MockUserRepositoryWithError) Delete(ctx context.Context, id int) error {
	return nil
}

func (m *MockUserRepositoryWithError) Exists(ctx context.Context, email string) (bool, error) {
	if m.existsError {
		return false, domain.ErrUserNotFound
	}
	return false, nil
}
