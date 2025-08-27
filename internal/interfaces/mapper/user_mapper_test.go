package mapper

import (
	"testing"
	"time"

	"hello-world/internal/domain"
	"hello-world/internal/interfaces/dto"
)

func TestUserMapper_ToUserResponse(t *testing.T) {
	mapper := NewUserMapper()

	user := &domain.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	response := mapper.ToUserResponse(user)

	if response.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, response.ID)
	}
	if response.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, response.Email)
	}
	if response.FirstName != user.FirstName {
		t.Errorf("Expected first name %s, got %s", user.FirstName, response.FirstName)
	}
	if response.LastName != user.LastName {
		t.Errorf("Expected last name %s, got %s", user.LastName, response.LastName)
	}
	if response.Phone != user.Phone {
		t.Errorf("Expected phone %s, got %s", user.Phone, response.Phone)
	}
	if !response.Birthday.Equal(user.Birthday) {
		t.Errorf("Expected birthday %v, got %v", user.Birthday, response.Birthday)
	}
	if !response.CreatedAt.Equal(user.CreatedAt) {
		t.Errorf("Expected created at %v, got %v", user.CreatedAt, response.CreatedAt)
	}
	if !response.UpdatedAt.Equal(user.UpdatedAt) {
		t.Errorf("Expected updated at %v, got %v", user.UpdatedAt, response.UpdatedAt)
	}
}

func TestUserMapper_ToLoginResponse(t *testing.T) {
	mapper := NewUserMapper()

	token := "test_token"
	user := &domain.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Birthday:  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	response := mapper.ToLoginResponse(token, user)

	if response.Token != token {
		t.Errorf("Expected token %s, got %s", token, response.Token)
	}
	if response.User.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, response.User.ID)
	}
	if response.User.Email != user.Email {
		t.Errorf("Expected user email %s, got %s", user.Email, response.User.Email)
	}
}

func TestUserMapper_ParseCreateUserRequest(t *testing.T) {
	mapper := NewUserMapper()

	tests := []struct {
		name        string
		req         dto.CreateUserRequest
		expectError bool
		expectedErr error
	}{
		{
			name: "Valid request",
			req: dto.CreateUserRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "1234567890",
				Birthday:  "1990-01-01",
			},
			expectError: false,
		},
		{
			name: "Invalid birthday format",
			req: dto.CreateUserRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "1234567890",
				Birthday:  "invalid-date",
			},
			expectError: true,
			expectedErr: domain.ErrInvalidBirthday,
		},
		{
			name: "Empty birthday",
			req: dto.CreateUserRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "1234567890",
				Birthday:  "",
			},
			expectError: true,
			expectedErr: domain.ErrInvalidBirthday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, password, firstName, lastName, phone, birthday, err := mapper.ParseCreateUserRequest(tt.req)

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

			if email != tt.req.Email {
				t.Errorf("Expected email %s, got %s", tt.req.Email, email)
			}
			if password != tt.req.Password {
				t.Errorf("Expected password %s, got %s", tt.req.Password, password)
			}
			if firstName != tt.req.FirstName {
				t.Errorf("Expected first name %s, got %s", tt.req.FirstName, firstName)
			}
			if lastName != tt.req.LastName {
				t.Errorf("Expected last name %s, got %s", tt.req.LastName, lastName)
			}
			if phone != tt.req.Phone {
				t.Errorf("Expected phone %s, got %s", tt.req.Phone, phone)
			}

			expectedBirthday, _ := time.Parse("2006-01-02", tt.req.Birthday)
			if !birthday.Equal(expectedBirthday) {
				t.Errorf("Expected birthday %v, got %v", expectedBirthday, birthday)
			}
		})
	}
}

func TestUserMapper_ParseUpdateUserRequest(t *testing.T) {
	mapper := NewUserMapper()

	tests := []struct {
		name        string
		req         dto.UpdateUserRequest
		expectError bool
		expectedErr error
	}{
		{
			name: "Valid request with all fields",
			req: dto.UpdateUserRequest{
				FirstName: "Jane",
				LastName:  "Smith",
				Phone:     "9876543210",
				Birthday:  "1992-05-15",
			},
			expectError: false,
		},
		{
			name: "Valid request with partial fields",
			req: dto.UpdateUserRequest{
				FirstName: "Jane",
				LastName:  "Smith",
			},
			expectError: false,
		},
		{
			name:        "Valid request with empty fields",
			req:         dto.UpdateUserRequest{},
			expectError: false,
		},
		{
			name: "Invalid birthday format",
			req: dto.UpdateUserRequest{
				FirstName: "Jane",
				LastName:  "Smith",
				Phone:     "9876543210",
				Birthday:  "invalid-date",
			},
			expectError: true,
			expectedErr: domain.ErrInvalidBirthday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firstName, lastName, phone, birthday, err := mapper.ParseUpdateUserRequest(tt.req)

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

			if firstName != tt.req.FirstName {
				t.Errorf("Expected first name %s, got %s", tt.req.FirstName, firstName)
			}
			if lastName != tt.req.LastName {
				t.Errorf("Expected last name %s, got %s", tt.req.LastName, lastName)
			}
			if phone != tt.req.Phone {
				t.Errorf("Expected phone %s, got %s", tt.req.Phone, phone)
			}

			if tt.req.Birthday != "" {
				if birthday == nil {
					t.Errorf("Expected birthday to be set but got nil")
				} else {
					expectedBirthday, _ := time.Parse("2006-01-02", tt.req.Birthday)
					if !birthday.Equal(expectedBirthday) {
						t.Errorf("Expected birthday %v, got %v", expectedBirthday, *birthday)
					}
				}
			} else {
				if birthday != nil {
					t.Errorf("Expected birthday to be nil but got %v", *birthday)
				}
			}
		})
	}
}

func TestNewUserMapper(t *testing.T) {
	mapper := NewUserMapper()
	if mapper == nil {
		t.Error("Expected mapper to be created but got nil")
	}
}
