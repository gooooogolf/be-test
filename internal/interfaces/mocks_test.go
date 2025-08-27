package interfaces

import (
	"context"
	"time"

	"hello-world/internal/domain"
)

// MockUserService for testing
type MockUserService struct{}

func (m *MockUserService) Register(ctx context.Context, email, password, firstName, lastName, phone string, birthday time.Time) (*domain.User, error) {
	return &domain.User{ID: 1, Email: email}, nil
}

func (m *MockUserService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	return "test_token", &domain.User{ID: 1, Email: email}, nil
}

func (m *MockUserService) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	return &domain.User{ID: userID, Email: "test@example.com"}, nil
}

func (m *MockUserService) GetUserProfile(ctx context.Context, userID int) (*domain.User, error) {
	return &domain.User{ID: userID, Email: "test@example.com"}, nil
}

func (m *MockUserService) UpdateUser(ctx context.Context, userID int, firstName, lastName, phone string, birthday *time.Time) (*domain.User, error) {
	return &domain.User{ID: userID, FirstName: firstName, LastName: lastName}, nil
}

// MockAuthServiceForRouter for testing router specifically (different from middleware mock)
type MockAuthServiceForRouter struct{}

func (m *MockAuthServiceForRouter) HashPassword(password string) (string, error) {
	return "hashed_" + password, nil
}

func (m *MockAuthServiceForRouter) ComparePassword(hashedPassword, password string) error {
	return nil
}

func (m *MockAuthServiceForRouter) GenerateToken(userID int, email string) (string, error) {
	return "test_token", nil
}

func (m *MockAuthServiceForRouter) ValidateToken(token string) (*domain.TokenClaims, error) {
	return &domain.TokenClaims{UserID: 1, Email: "test@example.com"}, nil
}
