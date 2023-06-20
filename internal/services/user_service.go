package services

import (
	"strings"
	"virtual-file-system/internal/models"
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
	name = strings.ToLower(name)
	s.Users[name] = models.User{Name: name}
	return nil
}

func (s *UserService) Exist(name string) bool {
	name = strings.ToLower(name)
	_, exist := s.Users[name]
	return exist
}
