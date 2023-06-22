package services

import (
	"testing"
	"virtual-file-system/internal/models"
)

func TestFileService_Creation(t *testing.T) {
	testCases := []struct {
		name         string
		users        map[string]models.User
		targetUser   string
		targetFolder string
		targetFile   string
		description  string
		expectedErr  string
		expectedLen  int
		expectedName string
	}{
		{
			name: "Create a file in an existing folder",
			users: map[string]models.User{"dalaoqi": {Name: "dalaoqi", Folders: map[string]models.Folder{"myfolder": {
				Name:        "myfolder",
				Description: "My folder description",
				Files:       nil,
			}}}},
			targetUser:   "dalaoqi",
			targetFolder: "myfolder",
			targetFile:   "myfile",
			description:  "My file description",
			expectedErr:  "",
			expectedLen:  1,
			expectedName: "myfile",
		},
		{
			name: "Create a file in a non-existing folder",
			users: map[string]models.User{"dalaoqi": {Name: "dalaoqi", Folders: map[string]models.Folder{"myfolder": {
				Name:        "myfolder",
				Description: "My folder description",
				Files:       nil,
			}}}},
			targetUser:   "dalaoqi",
			targetFolder: "otherfolder",
			targetFile:   "myfile",
			description:  "My file description",
			expectedErr:  "Error: The otherfolder doesn't exist.",
			expectedLen:  0,
			expectedName: "",
		},
		{
			name: "Create a file with duplicated name for the folder",
			users: map[string]models.User{"dalaoqi": {Name: "dalaoqi", Folders: map[string]models.Folder{"myfolder": {
				Name:        "myfolder",
				Description: "My folder description",
				Files:       nil,
			}}}},
			targetUser:   "dalaoqi",
			targetFolder: "myfolder",
			targetFile:   "myfile",
			description:  "My file description",
			expectedErr:  "Error: The myfile has already existed.",
			expectedLen:  1,
			expectedName: "myfile",
		},
		{
			name: "Create a folder with invalid chars",
			users: map[string]models.User{"dalaoqi": {Name: "dalaoqi", Folders: map[string]models.Folder{"myfolder": {
				Name:        "myfolder",
				Description: "My folder description",
				Files:       nil,
			}}}},
			targetUser:   "dalaoqi",
			targetFolder: "myfolder",
			targetFile:   "myfile^%$",
			description:  "My file description",
			expectedErr:  "Error: The myfile^%$ contains invalid chars.",
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
			fileService := NewFileService(userService, folderService)

			// Perform the test by calling fileService.CreateFolder() and check the error message
			if err := fileService.CreateFile(test.targetUser, test.targetFolder, test.targetFile, test.description); err != nil && err.Error() != test.expectedErr {
				t.Errorf("fileService.CreateFile() has error: %v, expected: %v", err.Error(), test.expectedErr)
			}

			if len(fileService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files) != test.expectedLen {
				t.Errorf("fileService.CreateFile() len = %v, expectedLen %v", len(folderService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files), test.expectedLen)
			}

			if test.expectedName != "" && fileService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files[test.targetFile].Name != test.expectedName {
				t.Errorf("fileService.CreateFile() name = %v, expectedName %v", folderService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files[test.targetFile].Name, test.expectedName)
			}

			// Verify the created file's attributes
			if len(folderService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files) > 0 {
				file := folderService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files[test.targetFile]
				if file.Description != test.description {
					t.Errorf("fileService.CreateFolder() description = %v, expectedDescription %v", file.Description, test.description)
				}
			}
		})
	}
}
