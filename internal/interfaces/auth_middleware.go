package interfaces

import (
	"context"
	"net/http"
	"strings"

	"hello-world/internal/domain"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	authService domain.AuthService
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(authService domain.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Middleware returns the HTTP middleware function
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.sendUnauthorizedResponse(w, "Authorization header required")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			m.sendUnauthorizedResponse(w, "Bearer token required")
			return
		}

		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			m.sendUnauthorizedResponse(w, "Invalid token")
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) sendUnauthorizedResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
