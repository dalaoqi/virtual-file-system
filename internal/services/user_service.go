package services

import (
	"fmt"
	"virtual-file-system/internal/models"
	"virtual-file-system/internal/utils"
)

// UserService handles user-related operations
type UserService struct {
	Users map[string]models.User
}

// NewUserService creates a new instance of UserService
func NewUserService() *UserService {
	return &UserService{
		Users: make(map[string]models.User),
	}
}

// Register registers a new user
func (s *UserService) Register(name string) error {
	// Check if the user already exists
	if _, exist := s.Users[name]; exist {
		return fmt.Errorf("Error: The %v has already existed.", name)
	}

	// Check if the name contains invalid characters
	if utils.ExistInvalidChars(name) {
		return fmt.Errorf("Error: The %v contains invalid chars.", name)
	}

	// Register the new user
	s.Users[name] = models.User{Name: name}
	return nil
}
