package services

import (
	"fmt"
	"strings"
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
func (s *UserService) Register(userName string) error {
	userName = strings.ToLower(userName)
	// Check if the user already exists
	if s.Exist(userName) {
		return fmt.Errorf("Error: The %v has already existed.", userName)
	}

	// Check if the name contains invalid characters
	if utils.ExistInvalidChars(userName) {
		return fmt.Errorf("Error: The %v contains invalid chars.", userName)
	}

	s.Users[userName] = models.User{Name: userName, Folders: nil}
	return nil
}

func (s *UserService) Exist(name string) bool {
	_, exist := s.Users[name]
	return exist
}
