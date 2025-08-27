package infrastructure

import (
	"testing"

	"hello-world/internal/domain"
)

func TestNewJWTAuthService(t *testing.T) {
	authService := NewJWTAuthService()
	if authService == nil {
		t.Error("Expected auth service to be created but got nil")
	}
}

func TestJWTAuthService_HashPassword(t *testing.T) {
	authService := NewJWTAuthService()
	password := "testpassword123"

	hashedPassword, err := authService.HashPassword(password)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if hashedPassword == "" {
		t.Error("Expected hashed password to be non-empty")
	}

	if hashedPassword == password {
		t.Error("Expected hashed password to be different from original password")
	}
}

func TestJWTAuthService_ComparePassword(t *testing.T) {
	authService := NewJWTAuthService()
	password := "testpassword123"

	hashedPassword, err := authService.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	err = authService.ComparePassword(hashedPassword, password)
	if err != nil {
		t.Errorf("Expected no error for correct password, got: %v", err)
	}

	// Test incorrect password
	err = authService.ComparePassword(hashedPassword, "wrongpassword")
	if err == nil {
		t.Error("Expected error for incorrect password")
	}
}

func TestJWTAuthService_GenerateToken(t *testing.T) {
	authService := NewJWTAuthService()
	userID := 123
	email := "test@example.com"

	token, err := authService.GenerateToken(userID, email)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if token == "" {
		t.Error("Expected token to be non-empty")
	}

	// Token should contain the JWT structure (header.payload.signature)
	tokenParts := len(token)
	if tokenParts < 10 { // JWT tokens are much longer
		t.Error("Generated token seems too short to be a valid JWT")
	}
}

func TestJWTAuthService_ValidateToken(t *testing.T) {
	authService := NewJWTAuthService()
	userID := 123
	email := "test@example.com"

	// Generate a token first
	token, err := authService.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	claims, err := authService.ValidateToken(token)
	if err != nil {
		t.Errorf("Unexpected error validating token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}
}

func TestJWTAuthService_ValidateToken_Invalid(t *testing.T) {
	authService := NewJWTAuthService()

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "Empty token",
			token: "",
		},
		{
			name:  "Invalid token format",
			token: "invalid.token.format",
		},
		{
			name:  "Malformed token",
			token: "not-a-jwt-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := authService.ValidateToken(tt.token)
			if err == nil {
				t.Error("Expected error for invalid token")
			}
		})
	}
}

func TestJWTAuthService_TokenRoundTrip(t *testing.T) {
	authService := NewJWTAuthService()

	testCases := []struct {
		userID int
		email  string
	}{
		{1, "user1@example.com"},
		{999, "user999@example.com"},
		{42, "test@domain.co.uk"},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			// Generate token
			token, err := authService.GenerateToken(tc.userID, tc.email)
			if err != nil {
				t.Fatalf("Failed to generate token: %v", err)
			}

			// Validate token
			claims, err := authService.ValidateToken(token)
			if err != nil {
				t.Fatalf("Failed to validate token: %v", err)
			}

			// Check claims
			if claims.UserID != tc.userID {
				t.Errorf("Expected user ID %d, got %d", tc.userID, claims.UserID)
			}
			if claims.Email != tc.email {
				t.Errorf("Expected email %s, got %s", tc.email, claims.Email)
			}
		})
	}
}

func TestJWTAuthService_Interface(t *testing.T) {
	// Test that JWTAuthService implements domain.AuthService interface
	var _ domain.AuthService = &JWTAuthService{}
}
