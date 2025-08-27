package interfaces

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	mockUserService := &MockUserService{}
	mockAuthService := &MockAuthServiceForRouter{}

	router := NewRouter(mockUserService, mockAuthService)

	if router == nil {
		t.Error("Expected router to be created")
	}

	if router.userHandler == nil {
		t.Error("Expected user handler to be set")
	}

	if router.authMiddleware == nil {
		t.Error("Expected auth middleware to be set")
	}
}

func TestRouter_SetupRoutes(t *testing.T) {
	mockUserService := &MockUserService{}
	mockAuthService := &MockAuthServiceForRouter{}
	router := NewRouter(mockUserService, mockAuthService)

	chiRouter := router.SetupRoutes()

	if chiRouter == nil {
		t.Error("Expected chi router to be created")
	}
}

func TestRouter_PublicRoutes(t *testing.T) {
	mockUserService := &MockUserService{}
	mockAuthService := &MockAuthServiceForRouter{}
	router := NewRouter(mockUserService, mockAuthService)
	chiRouter := router.SetupRoutes()

	// Test cases for public routes
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "Hello endpoint",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "Register endpoint accepts POST",
			method:         "POST",
			path:           "/register",
			expectedStatus: http.StatusBadRequest, // Will fail validation but route exists
			checkResponse:  false,
		},
		{
			name:           "Login endpoint accepts POST",
			method:         "POST",
			path:           "/login",
			expectedStatus: http.StatusBadRequest, // Will fail validation but route exists
			checkResponse:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rr := httptest.NewRecorder()

			chiRouter.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.checkResponse && rr.Body.String() == "" {
				t.Error("Expected response body but got empty")
			}
		})
	}
}

func TestRouter_ProtectedRoutes(t *testing.T) {
	mockUserService := &MockUserService{}
	mockAuthService := &MockAuthServiceForRouter{}
	router := NewRouter(mockUserService, mockAuthService)
	chiRouter := router.SetupRoutes()

	// Test protected route with valid token
	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	rr := httptest.NewRecorder()

	chiRouter.ServeHTTP(rr, req)

	// Should return 200 since our mock always validates successfully
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestRouter_ProtectedRoutes_NoAuth(t *testing.T) {
	mockUserService := &MockUserService{}
	mockAuthService := &MockAuthServiceForRouter{}
	router := NewRouter(mockUserService, mockAuthService)
	chiRouter := router.SetupRoutes()

	// Test protected route without token
	req := httptest.NewRequest("GET", "/me", nil)
	rr := httptest.NewRecorder()

	chiRouter.ServeHTTP(rr, req)

	// Should return 401 unauthorized
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", rr.Code)
	}
}

func TestRouter_NotFoundRoute(t *testing.T) {
	mockUserService := &MockUserService{}
	mockAuthService := &MockAuthServiceForRouter{}
	router := NewRouter(mockUserService, mockAuthService)
	chiRouter := router.SetupRoutes()

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	rr := httptest.NewRecorder()

	chiRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	mockUserService := &MockUserService{}
	mockAuthService := &MockAuthServiceForRouter{}
	router := NewRouter(mockUserService, mockAuthService)
	chiRouter := router.SetupRoutes()

	// Try to POST to hello endpoint which only accepts GET
	req := httptest.NewRequest("POST", "/", nil)
	rr := httptest.NewRecorder()

	chiRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

func TestRouter_SwaggerEndpoint(t *testing.T) {
	mockUserService := &MockUserService{}
	mockAuthService := &MockAuthServiceForRouter{}
	router := NewRouter(mockUserService, mockAuthService)
	chiRouter := router.SetupRoutes()

	req := httptest.NewRequest("GET", "/swagger/", nil)
	rr := httptest.NewRecorder()

	chiRouter.ServeHTTP(rr, req)

	// Swagger endpoint should be accessible
	if rr.Code == http.StatusNotFound {
		t.Error("Swagger endpoint should be accessible")
	}
}
