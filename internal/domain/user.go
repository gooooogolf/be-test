package domain

import (
	"context"
	"time"
)

// User represents the core user entity in the domain
type User struct {
	ID        int
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
	Birthday  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user entity with validation
func NewUser(email, password, firstName, lastName, phone string, birthday time.Time) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if firstName == "" {
		return nil, ErrInvalidFirstName
	}
	if lastName == "" {
		return nil, ErrInvalidLastName
	}

	now := time.Now()
	return &User{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		Birthday:  birthday,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsValidForUpdate checks if user data is valid for update
func (u *User) IsValidForUpdate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.FirstName == "" {
		return ErrInvalidFirstName
	}
	if u.LastName == "" {
		return ErrInvalidLastName
	}
	return nil
}

// UserRepository defines the contract for user data persistence
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, email string) (bool, error)
}

// AuthService defines the contract for authentication operations
type AuthService interface {
	GenerateToken(userID int, email string) (string, error)
	ValidateToken(token string) (*TokenClaims, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
}

// UserService defines the use case interface for user operations
type UserService interface {
	Register(ctx context.Context, email, password, firstName, lastName, phone string, birthday time.Time) (*User, error)
	Login(ctx context.Context, email, password string) (string, *User, error)
	GetUserByID(ctx context.Context, userID int) (*User, error)
	GetUserProfile(ctx context.Context, userID int) (*User, error)
	UpdateUser(ctx context.Context, userID int, firstName, lastName, phone string, birthday *time.Time) (*User, error)
}

// DomainError represents domain-specific errors
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e DomainError) Error() string {
	return e.Message
}

// Domain errors
var (
	ErrUserNotFound         = DomainError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrUserAlreadyExists    = DomainError{Code: "USER_ALREADY_EXISTS", Message: "User already exists"}
	ErrInvalidCredentials   = DomainError{Code: "INVALID_CREDENTIALS", Message: "Invalid email or password"}
	ErrInvalidToken         = DomainError{Code: "INVALID_TOKEN", Message: "Invalid token"}
	ErrUnauthorized         = DomainError{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
	ErrInvalidEmail         = DomainError{Code: "INVALID_EMAIL", Message: "Email is required"}
	ErrInvalidFirstName     = DomainError{Code: "INVALID_FIRST_NAME", Message: "First name is required"}
	ErrInvalidLastName      = DomainError{Code: "INVALID_LAST_NAME", Message: "Last name is required"}
	ErrInvalidBirthday      = DomainError{Code: "INVALID_BIRTHDAY", Message: "Invalid birthday format (YYYY-MM-DD)"}
	ErrPasswordHashError    = DomainError{Code: "PASSWORD_HASH_ERROR", Message: "Failed to hash password"}
	ErrUserCreationError    = DomainError{Code: "USER_CREATION_ERROR", Message: "Failed to create user"}
	ErrTokenGenerationError = DomainError{Code: "TOKEN_GENERATION_ERROR", Message: "Failed to generate token"}
)
