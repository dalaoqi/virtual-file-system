package services

import (
	"testing"
	"virtual-file-system/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestFolderService_Creation(t *testing.T) {
	testCases := []struct {
		name         string
		folders      map[string]map[string]models.Folder
		targetUser   string
		targetFolder string
		description  string
		expected     string
	}{
		{
			name:         "Create a new folder for the user",
			folders:      map[string]map[string]models.Folder{},
			targetUser:   "dalaoqi",
			targetFolder: "myfolder",
			description:  "My folder description",
			expected:     "",
		},
		{
			name: "Create a folder with duplicated name for the user",
			folders: map[string]map[string]models.Folder{
				"dalaoqi": {
					"myfolder": {
						Name:        "myfolder",
						Owner:       "dalaoqi",
						Description: "My folder description",
					},
				},
			},
			targetUser:   "dalaoqi",
			targetFolder: "myfolder",
			description:  "Another folder description",
			expected:     "Error: The folder name 'myfolder' already exists for the user 'dalaoqi'.",
		},
		{
			name:         "Create a folder with invalid chars",
			folders:      map[string]map[string]models.Folder{},
			targetUser:   "dalaoqi",
			targetFolder: "myfolder$%",
			description:  "Folder with invalid chars",
			expected:     "Error: The folder name 'myfolder$%' contains invalid chars.",
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			FolderService := &FolderService{
				Folders: test.folders,
			}
			// Perform the test by calling FolderService.CreateFolder() and check the error message
			if err := FolderService.CreateFolder(test.targetUser, test.targetFolder, test.description); err != nil && err.Error() != test.expected {
				t.Errorf("FolderService.CreateFolder() has error: %v, expected: %v", err.Error(), test.expected)
			}
		})
	}
}

func TestFolderService_Exist(t *testing.T) {
	service := NewFolderService()
	service.Folders["dalaoqi"] = map[string]models.Folder{
		"folder1": {Owner: "dalaoqi", Name: "folder1"},
		"folder2": {Owner: "dalaoqi", Name: "folder2"},
	}

	testCases := []struct {
		name       string
		username   string
		foldername string
		expected   bool
	}{
		{
			name:       "Existing folder name",
			username:   "dalaoqi",
			foldername: "folder1",
			expected:   true,
		},
		{
			name:       "Existing folder name (upper)",
			username:   "DALAOQI",
			foldername: "FOLDER1",
			expected:   true,
		},
		{
			name:       "Non-existing folder name",
			username:   "dalaoqi",
			foldername: "folder3",
			expected:   false,
		},
		{
			name:       "Non-existing user",
			username:   "dalaoqi123",
			foldername: "folder1",
			expected:   false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			exist := service.Exist(test.username, test.foldername)
			assert.Equal(t, test.expected, exist)
		})
	}
}
