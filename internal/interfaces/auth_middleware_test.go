package interfaces

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"hello-world/internal/domain"
)

// MockAuthService for testing middleware
type MockAuthService struct {
	ValidateTokenFunc func(token string) (*domain.TokenClaims, error)
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

func (m *MockAuthService) GenerateToken(userID int, email string) (string, error) {
	return "valid_token", nil
}

func (m *MockAuthService) ValidateToken(token string) (*domain.TokenClaims, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(token)
	}
	return nil, domain.ErrInvalidToken
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Setup mock auth service
	mockAuth := &MockAuthService{
		ValidateTokenFunc: func(token string) (*domain.TokenClaims, error) {
			if token == "valid_token" {
				return &domain.TokenClaims{
					UserID: 1,
					Email:  "test@example.com",
				}, nil
			}
			return nil, domain.ErrInvalidToken
		},
	}

	middleware := NewAuthMiddleware(mockAuth)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(int)
		if !ok || userID != 1 {
			t.Error("Expected userID 1 in context")
		}
		email, ok := r.Context().Value("email").(string)
		if !ok || email != "test@example.com" {
			t.Error("Expected email test@example.com in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with middleware
	handler := middleware.Middleware(testHandler)

	// Create request with valid token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	rr := httptest.NewRecorder()

	// Execute request
	handler.ServeHTTP(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	mockAuth := &MockAuthService{}
	middleware := NewAuthMiddleware(mockAuth)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	handler := middleware.Middleware(testHandler)

	// Create request without token
	req := httptest.NewRequest("GET", "/protected", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check response
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockAuth := &MockAuthService{
		ValidateTokenFunc: func(token string) (*domain.TokenClaims, error) {
			return nil, domain.ErrInvalidToken
		},
	}

	middleware := NewAuthMiddleware(mockAuth)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	handler := middleware.Middleware(testHandler)

	// Create request with invalid token
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check response
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestAuthMiddleware_InvalidTokenFormat(t *testing.T) {
	mockAuth := &MockAuthService{}
	middleware := NewAuthMiddleware(mockAuth)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	handler := middleware.Middleware(testHandler)

	// Test cases for invalid token formats
	testCases := []struct {
		name           string
		authorization  string
		expectedStatus int
	}{
		{
			name:           "Token without Bearer prefix",
			authorization:  "invalid_token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Bearer without token",
			authorization:  "Bearer",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty Bearer token",
			authorization:  "Bearer ",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", tc.authorization)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestNewAuthMiddleware(t *testing.T) {
	mockAuth := &MockAuthService{}
	middleware := NewAuthMiddleware(mockAuth)

	if middleware == nil {
		t.Error("Expected middleware to be created")
	}

	if middleware.authService != mockAuth {
		t.Error("Expected auth service to be set correctly")
	}
}
