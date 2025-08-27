package domain

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		firstName   string
		lastName    string
		phone       string
		birthday    time.Time
		expectError bool
		expectedErr error
	}{
		{
			name:        "Valid user creation",
			email:       "test@example.com",
			password:    "hashedpassword",
			firstName:   "John",
			lastName:    "Doe",
			phone:       "1234567890",
			birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "Empty email",
			email:       "",
			password:    "hashedpassword",
			firstName:   "John",
			lastName:    "Doe",
			phone:       "1234567890",
			birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: true,
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "Empty first name",
			email:       "test@example.com",
			password:    "hashedpassword",
			firstName:   "",
			lastName:    "Doe",
			phone:       "1234567890",
			birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: true,
			expectedErr: ErrInvalidFirstName,
		},
		{
			name:        "Empty last name",
			email:       "test@example.com",
			password:    "hashedpassword",
			firstName:   "John",
			lastName:    "",
			phone:       "1234567890",
			birthday:    time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: true,
			expectedErr: ErrInvalidLastName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.password, tt.firstName, tt.lastName, tt.phone, tt.birthday)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if err != tt.expectedErr {
					t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user.Email != tt.email {
				t.Errorf("Expected email %s, got %s", tt.email, user.Email)
			}
			if user.Password != tt.password {
				t.Errorf("Expected password %s, got %s", tt.password, user.Password)
			}
			if user.FirstName != tt.firstName {
				t.Errorf("Expected first name %s, got %s", tt.firstName, user.FirstName)
			}
			if user.LastName != tt.lastName {
				t.Errorf("Expected last name %s, got %s", tt.lastName, user.LastName)
			}
			if user.Phone != tt.phone {
				t.Errorf("Expected phone %s, got %s", tt.phone, user.Phone)
			}
			if !user.Birthday.Equal(tt.birthday) {
				t.Errorf("Expected birthday %v, got %v", tt.birthday, user.Birthday)
			}
		})
	}
}

func TestUser_GetFullName(t *testing.T) {
	user := &User{
		FirstName: "John",
		LastName:  "Doe",
	}

	expected := "John Doe"
	actual := user.GetFullName()

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestUser_IsValidForUpdate(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		expectError bool
		expectedErr error
	}{
		{
			name: "Valid user for update",
			user: &User{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectError: false,
		},
		{
			name: "Invalid email for update",
			user: &User{
				Email:     "",
				FirstName: "John",
				LastName:  "Doe",
			},
			expectError: true,
			expectedErr: ErrInvalidEmail,
		},
		{
			name: "Invalid first name for update",
			user: &User{
				Email:     "test@example.com",
				FirstName: "",
				LastName:  "Doe",
			},
			expectError: true,
			expectedErr: ErrInvalidFirstName,
		},
		{
			name: "Invalid last name for update",
			user: &User{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "",
			},
			expectError: true,
			expectedErr: ErrInvalidLastName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.IsValidForUpdate()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if err != tt.expectedErr {
					t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDomainError_Error(t *testing.T) {
	err := DomainError{
		Code:    "TEST_ERROR",
		Message: "Test error message",
	}

	expected := "Test error message"
	actual := err.Error()

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestDomainErrors(t *testing.T) {
	// Test all predefined domain errors
	errors := []DomainError{
		ErrUserNotFound,
		ErrUserAlreadyExists,
		ErrInvalidCredentials,
		ErrInvalidToken,
		ErrUnauthorized,
		ErrInvalidEmail,
		ErrInvalidFirstName,
		ErrInvalidLastName,
		ErrInvalidBirthday,
		ErrPasswordHashError,
		ErrUserCreationError,
		ErrTokenGenerationError,
	}

	for _, err := range errors {
		if err.Code == "" {
			t.Errorf("Error code should not be empty for error: %v", err)
		}
		if err.Message == "" {
			t.Errorf("Error message should not be empty for error: %v", err)
		}
		if err.Error() != err.Message {
			t.Errorf("Error() should return message for error: %v", err)
		}
	}
}
