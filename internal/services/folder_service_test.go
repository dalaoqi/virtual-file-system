package services

import (
	"reflect"
	"testing"
	"time"
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

func TestFolderService_GetFolders(t *testing.T) {
	folderService := NewFolderService()

	// Create some folders for a user
	folder1 := models.Folder{
		Name:        "folder1",
		Owner:       "dalaoqi",
		Description: "Folder 1",
		CreatedAt:   time.Now().Add(-time.Hour),
	}
	folder2 := models.Folder{
		Name:        "folder2",
		Owner:       "dalaoqi",
		Description: "Folder 2",
		CreatedAt:   time.Now().Add(-2 * time.Hour),
	}
	folder3 := models.Folder{
		Name:        "folder3",
		Owner:       "dalaoqi",
		Description: "Folder 3",
		CreatedAt:   time.Now().Add(-3 * time.Hour),
	}

	foldersMap := make(map[string]models.Folder)
	foldersMap[folder1.Name] = folder1
	foldersMap[folder2.Name] = folder2
	foldersMap[folder3.Name] = folder3

	folderService.Folders["dalaoqi"] = foldersMap
	folderService.Folders["dalaoqiEmpty"] = make(map[string]models.Folder)

	testCases := []struct {
		name       string
		userName   string
		sortField  string
		sortOrder  string
		wantResult []models.Folder
	}{
		{
			name:       "Sort by name in ascending order",
			userName:   "dalaoqi",
			sortField:  "name",
			sortOrder:  "asc",
			wantResult: []models.Folder{folder1, folder2, folder3},
		},
		{
			name:       "Sort by name in descending order",
			userName:   "dalaoqi",
			sortField:  "name",
			sortOrder:  "desc",
			wantResult: []models.Folder{folder3, folder2, folder1},
		},
		{
			name:       "Sort by created at in ascending order",
			userName:   "dalaoqi",
			sortField:  "created",
			sortOrder:  "asc",
			wantResult: []models.Folder{folder3, folder2, folder1},
		},
		{
			name:       "Sort by created at in descending order",
			userName:   "dalaoqi",
			sortField:  "created",
			sortOrder:  "desc",
			wantResult: []models.Folder{folder1, folder2, folder3},
		},
		{
			name:       "User doesn't have any folders",
			userName:   "dalaoqiEmpty",
			sortField:  "name",
			sortOrder:  "asc",
			wantResult: []models.Folder{},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			gotResult, err := folderService.GetFolders(test.userName, test.sortField, test.sortOrder)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(gotResult, test.wantResult) {
				t.Errorf("Result mismatch, Got: %v, Want: %v", gotResult, test.wantResult)
			}
		})
	}
}
