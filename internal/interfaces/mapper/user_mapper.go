package mapper

import (
	"time"

	"hello-world/internal/domain"
	"hello-world/internal/interfaces/dto"
)

// UserMapper handles conversion between domain entities and DTOs
type UserMapper struct{}

// NewUserMapper creates a new UserMapper instance
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ToUserResponse converts a domain User to a UserResponse DTO
func (m *UserMapper) ToUserResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Birthday:  user.Birthday,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToLoginResponse converts login data to a LoginResponse DTO
func (m *UserMapper) ToLoginResponse(token string, user *domain.User) dto.LoginResponse {
	return dto.LoginResponse{
		Token: token,
		User:  m.ToUserResponse(user),
	}
}

// ParseCreateUserRequest converts a CreateUserRequest DTO to domain parameters
func (m *UserMapper) ParseCreateUserRequest(req dto.CreateUserRequest) (email, password, firstName, lastName, phone string, birthday time.Time, err error) {
	birthday, err = time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return "", "", "", "", "", time.Time{}, domain.ErrInvalidBirthday
	}

	return req.Email, req.Password, req.FirstName, req.LastName, req.Phone, birthday, nil
}

// ParseUpdateUserRequest converts an UpdateUserRequest DTO to domain parameters
func (m *UserMapper) ParseUpdateUserRequest(req dto.UpdateUserRequest) (firstName, lastName, phone string, birthday *time.Time, err error) {
	firstName = req.FirstName
	lastName = req.LastName
	phone = req.Phone

	if req.Birthday != "" {
		parsedBirthday, parseErr := time.Parse("2006-01-02", req.Birthday)
		if parseErr != nil {
			return "", "", "", nil, domain.ErrInvalidBirthday
		}
		birthday = &parsedBirthday
	}

	return firstName, lastName, phone, birthday, nil
}
