package usecase

import (
	"context"
	"time"

	"hello-world/internal/domain"
)

// UserUseCase implements domain.UserService and handles user-related business logic
type UserUseCase struct {
	userRepo    domain.UserRepository
	authService domain.AuthService
}

// NewUserUseCase creates a new UserUseCase instance
func NewUserUseCase(userRepo domain.UserRepository, authService domain.AuthService) domain.UserService {
	return &UserUseCase{
		userRepo:    userRepo,
		authService: authService,
	}
}

// Register creates a new user account
func (uc *UserUseCase) Register(ctx context.Context, email, password, firstName, lastName, phone string, birthday time.Time) (*domain.User, error) {
	// Check if user already exists
	exists, err := uc.userRepo.Exists(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := uc.authService.HashPassword(password)
	if err != nil {
		return nil, domain.ErrPasswordHashError
	}

	// Create user entity with domain validation
	user, err := domain.NewUser(email, hashedPassword, firstName, lastName, phone, birthday)
	if err != nil {
		return nil, err
	}

	// Save user
	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, domain.ErrUserCreationError
	}

	// Create a copy for response to avoid modifying the stored user
	responseUser := *user
	responseUser.Password = ""
	return &responseUser, nil
}

// Login authenticates a user and returns a JWT token and user data
func (uc *UserUseCase) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, domain.ErrInvalidCredentials
	}

	// Verify password
	err = uc.authService.ComparePassword(user.Password, password)
	if err != nil {
		return "", nil, domain.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := uc.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", nil, domain.ErrTokenGenerationError
	}

	// Create a copy for response to avoid modifying the stored user
	responseUser := *user
	responseUser.Password = ""
	return token, &responseUser, nil
}

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Create a copy for response to avoid modifying the stored user
	responseUser := *user
	responseUser.Password = ""
	return &responseUser, nil
}

// GetUserProfile retrieves the current user's profile
func (uc *UserUseCase) GetUserProfile(ctx context.Context, userID int) (*domain.User, error) {
	return uc.GetUserByID(ctx, userID)
}

// UpdateUser updates user information
func (uc *UserUseCase) UpdateUser(ctx context.Context, userID int, firstName, lastName, phone string, birthday *time.Time) (*domain.User, error) {
	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Update fields if provided
	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}
	if phone != "" {
		user.Phone = phone
	}
	if birthday != nil {
		user.Birthday = *birthday
	}

	user.UpdatedAt = time.Now()

	// Validate updated user
	err = user.IsValidForUpdate()
	if err != nil {
		return nil, err
	}

	// Save updated user
	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, domain.ErrUserCreationError // Generic update error
	}

	// Create a copy for response to avoid modifying the stored user
	responseUser := *user
	responseUser.Password = ""
	return &responseUser, nil
}
