package services

import (
	"testing"
	"virtual-file-system/internal/models"
)

func TestUserRegister(t *testing.T) {
	testCases := []struct {
		name       string
		users      map[string]models.User
		targetName string
		expected   string
	}{
		{
			name:       "Add a new user to the empty users",
			users:      map[string]models.User{},
			targetName: "dalaoqi",
			expected:   "",
		},
		{
			name: "Add a duplicated user",
			users: map[string]models.User{
				"dalaoqi": {Name: "dalaoqi"},
			},
			targetName: "dalaoqi",
			expected:   "Error: The dalaoqi has already existed.",
		},
		{
			name: "Add a duplicated user (upper)",
			users: map[string]models.User{
				"dalaoqi": {Name: "dalaoqi"},
			},
			targetName: "DALAOQI",
			expected:   "Error: The dalaoqi has already existed.",
		},
		{
			name:       "Add a user with invalid chars",
			users:      map[string]models.User{},
			targetName: "dalaoqi&^%$?#.",
			expected:   "Error: The dalaoqi&^%$?#. contains invalid chars.",
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			UserService := &UserService{
				Users: test.users,
			}
			// Perform the test by calling UserService.Register() and check the error message
			if err := UserService.Register(test.targetName); err != nil && err.Error() != test.expected {
				t.Errorf("UserService.Register() has error: %s, expected: %s", err.Error(), test.expected)
			}
		})
	}
}
