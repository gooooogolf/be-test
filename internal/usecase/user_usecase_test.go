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
