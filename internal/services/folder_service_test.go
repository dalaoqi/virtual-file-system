package services

import (
	"reflect"
	"testing"
	"time"
	"virtual-file-system/internal/models"
)

func TestFolderService_Creation(t *testing.T) {
	testCases := []struct {
		name         string
		users        map[string]models.User
		targetUser   string
		targetFolder string
		description  string
		expectedErr  string
		expectedLen  int
		expectedName string
	}{
		{
			name:         "Create a new folder for the user",
			users:        map[string]models.User{"dalaoqi": {Name: "dalaoqi", Folders: nil}},
			targetUser:   "dalaoqi",
			targetFolder: "myfolder",
			description:  "My folder description",
			expectedErr:  "",
			expectedLen:  1,
			expectedName: "myfolder",
		},
		{
			name:         "Create a folder with duplicated name for the user",
			users:        map[string]models.User{"dalaoqi": {Name: "dalaoqi", Folders: map[string]models.Folder{"myfolder": {Name: "myfolder", Description: "My folder description"}}}},
			targetUser:   "dalaoqi",
			targetFolder: "myfolder",
			description:  "My folder description",
			expectedErr:  "Error: The myfolder has already existed.",
			expectedLen:  1,
			expectedName: "myfolder",
		},
		{
			name:         "Create a folder with invalid chars",
			users:        map[string]models.User{"dalaoqi": {Name: "dalaoqi", Folders: nil}},
			targetUser:   "dalaoqi",
			targetFolder: "myfolder^%$",
			description:  "My folder description",
			expectedErr:  "Error: The myfolder^%$ contains invalid chars.",
			expectedLen:  0,
			expectedName: "",
		},
		{
			name:         "Create a new folder for the non-exist user",
			users:        map[string]models.User{},
			targetUser:   "dalaoqi",
			targetFolder: "myfolder",
			description:  "My folder description",
			expectedErr:  "Error: The dalaoqi doesn't exist.",
			expectedLen:  0,
			expectedName: "",
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			userService := &UserService{
				Users: test.users,
			}

			folderService := NewFolderService(userService)

			// Perform the test by calling FolderService.CreateFolder() and check the error message
			if err := folderService.CreateFolder(test.targetUser, test.targetFolder, test.description); err != nil && err.Error() != test.expectedErr {
				t.Errorf("folderService.CreateFolder() has error: %v, expected: %v", err.Error(), test.expectedErr)
			}

			if len(folderService.UserService.Users[test.targetUser].Folders) != test.expectedLen {
				t.Errorf("folderService.CreateFolder() len = %v, expectedLen %v", len(folderService.UserService.Users[test.targetUser].Folders), test.expectedLen)
			}

			if test.expectedName != "" && folderService.UserService.Users[test.targetUser].Folders[test.targetFolder].Name != test.expectedName {
				t.Errorf("folderService.CreateFolder() name = %v, expectedName %v", folderService.UserService.Users[test.targetUser].Folders[test.targetFolder].Name, test.expectedName)
			}

			// Verify the created folder's attributes
			if len(folderService.UserService.Users[test.targetUser].Folders) > 0 {
				folder := folderService.UserService.Users[test.targetUser].Folders[test.targetFolder]
				if folder.Description != test.description {
					t.Errorf("folderService.CreateFolder() description = %v, expectedDescription %v", folder.Description, test.description)
				}
			}
		})
	}
}

func TestFolderService_GetFolders(t *testing.T) {
	// Create some folders for a user
	folder1 := models.Folder{
		Name:        "folder1",
		Description: "Folder 1",
		CreatedAt:   time.Now().Add(-time.Hour),
	}
	folder2 := models.Folder{
		Name:        "folder2",
		Description: "Folder 2",
		CreatedAt:   time.Now().Add(-2 * time.Hour),
	}
	folder3 := models.Folder{
		Name:        "folder3",
		Description: "Folder 3",
		CreatedAt:   time.Now().Add(-3 * time.Hour),
	}

	userService := &UserService{
		Users: map[string]models.User{
			"dalaoqi": {
				Name:    "dalaoqi",
				Folders: map[string]models.Folder{"folder1": folder1, "folder2": folder2, "folder3": folder3},
			},
			"dalaoqiEmpty": {
				Name:    "dalaoqiEmpty",
				Folders: nil,
			},
		},
	}

	folderService := NewFolderService(userService)

	testCases := []struct {
		name           string
		userName       string
		sortFlag       string
		sortOrderFlag  string
		expectedResult []models.Folder
		expectedError  string
	}{
		{
			name:           "Sort by name in ascending order",
			userName:       "dalaoqi",
			sortFlag:       "--sort-name",
			sortOrderFlag:  "asc",
			expectedResult: []models.Folder{folder1, folder2, folder3},
			expectedError:  "",
		},
		{
			name:           "Sort by name in descending order",
			userName:       "dalaoqi",
			sortFlag:       "--sort-name",
			sortOrderFlag:  "desc",
			expectedResult: []models.Folder{folder3, folder2, folder1},
			expectedError:  "",
		},
		{
			name:           "Sort by created at in ascending order",
			userName:       "dalaoqi",
			sortFlag:       "--sort-created",
			sortOrderFlag:  "asc",
			expectedResult: []models.Folder{folder3, folder2, folder1},
			expectedError:  "",
		},
		{
			name:           "Sort by created at in descending order",
			userName:       "dalaoqi",
			sortFlag:       "--sort-created",
			sortOrderFlag:  "desc",
			expectedResult: []models.Folder{folder1, folder2, folder3},
		},
		{
			name:           "Sort by an invalid flag",
			userName:       "testuser",
			sortFlag:       "--sort-invalid",
			sortOrderFlag:  "asc",
			expectedResult: []models.Folder{},
			expectedError:  "Error: Invalid sort flag",
		},
		{
			name:           "Sort by name in invalid order",
			userName:       "testuser",
			sortFlag:       "--sort-name",
			sortOrderFlag:  "invalid",
			expectedResult: []models.Folder{},
			expectedError:  "Error: Invalid sort order",
		},
		{
			name:           "User doesn't exist",
			userName:       "nonexistentuser",
			sortFlag:       "--sort-name",
			sortOrderFlag:  "asc",
			expectedResult: []models.Folder{},
			expectedError:  "Error: The nonexistentuser doesn't exist.",
		},
		{
			name:           "User doesn't have any folders",
			userName:       "dalaoqiEmpty",
			sortFlag:       "--sort-name",
			sortOrderFlag:  "asc",
			expectedResult: []models.Folder{},
			expectedError:  "Error: The dalaoqiEmpty doesn't exist.",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			gotResult, err := folderService.GetFolders(test.userName, test.sortFlag, test.sortOrderFlag)
			if err != nil && err.Error() != test.expectedError {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(gotResult, test.expectedResult) {
				t.Errorf("Result mismatch, Got: %v, Want: %v", gotResult, test.expectedResult)
			}
		})
	}
}

func TestFolderService_Deletion(t *testing.T) {
	testCases := []struct {
		name          string
		users         map[string]models.User
		targetUser    string
		targetFolder  string
		expectedError string
	}{
		{
			name: "Delete an existing folder for the user",
			users: map[string]models.User{
				"dalaoqi": {
					Name: "dalaoqi",
					Folders: map[string]models.Folder{"myfolder": {
						Name:        "myfolder",
						Description: "My folder description",
						CreatedAt:   time.Now(),
					},
					},
				},
			},
			targetUser:    "dalaoqi",
			targetFolder:  "myfolder",
			expectedError: "",
		},
		{
			name: "Delete a non-existing folder for the user",
			users: map[string]models.User{
				"dalaoqi": {
					Name: "dalaoqi",
					Folders: map[string]models.Folder{"myfolder": {
						Name:        "myfolder",
						Description: "My folder description",
						CreatedAt:   time.Now(),
					},
					},
				},
			},
			targetUser:    "dalaoqi",
			targetFolder:  "otherfolder",
			expectedError: "Error: The otherfolder doesn't exist",
		},
		{
			name: "Delete a folder for a non-existing user",
			users: map[string]models.User{
				"dalaoqi": {
					Name: "dalaoqi",
					Folders: map[string]models.Folder{"myfolder": {
						Name:        "myfolder",
						Description: "My folder description",
						CreatedAt:   time.Now(),
					},
					},
				},
			},
			targetUser:    "other",
			targetFolder:  "myfolder",
			expectedError: "Error: The other doesn't exist.",
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			userService := &UserService{
				Users: test.users,
			}
			folderService := &FolderService{
				UserService: userService,
			}
			// Perform the test by calling FolderService.DeleteFolder() and check the error message
			folderService.DeleteFolder(test.targetUser, test.targetFolder)

			// Check if the folder has been deleted from the folders map
			if exists := folderService.Exist(test.targetUser, test.targetFolder); exists {
				t.Errorf("Folder %s still exists for User %s", test.targetFolder, test.targetUser)
			}
		})
	}
}

func TestRenameFolder(t *testing.T) {
	testCases := []struct {
		name          string
		userName      string
		folderName    string
		newFolderName string
		expectedError string
	}{
		{
			name:          "Rename existing folder",
			userName:      "dalaoqi",
			folderName:    "folder1",
			newFolderName: "newfolder",
			expectedError: "",
		},
		{
			name:          "Rename non-existing folder",
			userName:      "dalaoqi",
			folderName:    "nonexistent",
			newFolderName: "newfolder",
			expectedError: "Error: The nonexistent doesn't exist",
		},
		{
			name:          "Rename folder with invalid characters",
			userName:      "dalaoqi",
			folderName:    "folder1",
			newFolderName: "new&folder",
			expectedError: "Error: The new&folder contains invalid chars.",
		},
		{
			name:          "Rename to an existing folder name",
			userName:      "dalaoqi",
			folderName:    "folder1",
			newFolderName: "folder2",
			expectedError: "Error: The folder2 has already existed.",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			userService := &UserService{
				Users: map[string]models.User{
					"dalaoqi": {
						Name: "dalaoqi",
						Folders: map[string]models.Folder{
							"folder1": {
								Name:        "folder1",
								Description: "Folder 1",
								CreatedAt:   time.Now(),
							},
							"folder2": {
								Name:        "folder2",
								Description: "Folder 2",
								CreatedAt:   time.Now(),
							},
						},
					},
				},
			}

			folderService := NewFolderService(userService)
			err := folderService.RenameFolder(test.userName, test.folderName, test.newFolderName)

			if test.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error: %s, but got nil", test.expectedError)
				} else if err.Error() != test.expectedError {
					t.Errorf("Expected error: %s, but got: %s", test.expectedError, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			}

			// Check if the folder was renamed correctly
			if err == nil {
				folders := userService.Users[test.userName].Folders
				_, exists := folders[test.folderName]
				if exists {
					t.Errorf("Folder %s should not exist", test.folderName)
				}

				renamedFolder, exists := folders[test.newFolderName]
				if !exists {
					t.Errorf("Renamed folder %s not found", test.newFolderName)
				} else if renamedFolder.Name != test.newFolderName {
					t.Errorf("Renamed folder name mismatch. Expected: %s, Got: %s", test.newFolderName, renamedFolder.Name)
				}
			}
		})
	}
}
