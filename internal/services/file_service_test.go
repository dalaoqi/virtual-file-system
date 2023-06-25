package services

import (
	"reflect"
	"testing"
	"time"
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
			name: "Create a file with invalid chars",
			users: map[string]models.User{"dalaoqi": {Name: "dalaoqi", Folders: map[string]models.Folder{"myfolder": {
				Name:        "myfolder",
				Description: "My folder description",
				Files:       nil,
			}}}},
			targetUser:   "dalaoqi",
			targetFolder: "myfolder",
			targetFile:   "myfile???",
			description:  "My file description",
			expectedErr:  "Error: The myfile??? contains invalid chars.",
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
				t.Errorf("fileService.CreateFile() has error: %s, expected: %s", err.Error(), test.expectedErr)
			}

			if len(fileService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files) != test.expectedLen {
				t.Errorf("fileService.CreateFile() len = %v, expectedLen %v", len(folderService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files), test.expectedLen)
			}

			if test.expectedName != "" && fileService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files[test.targetFile].Name != test.expectedName {
				t.Errorf("fileService.CreateFile() name = %s, expectedName %s", folderService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files[test.targetFile].Name, test.expectedName)
			}

			// Verify the created file's attributes
			if len(folderService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files) > 0 {
				file := folderService.UserService.Users[test.targetUser].Folders[test.targetFolder].Files[test.targetFile]
				if file.Description != test.description {
					t.Errorf("fileService.CreateFolder() description = %s, expectedDescription %s", file.Description, test.description)
				}
			}
		})
	}
}

func TestFileService_GetFiles(t *testing.T) {
	// Create some files for a user
	file1 := models.File{
		Name:      "myfile1",
		CreatedAt: time.Now().Add(-time.Hour),
	}
	file2 := models.File{
		Name:      "myfile2",
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}
	file3 := models.File{
		Name:      "myfile3",
		CreatedAt: time.Now().Add(-3 * time.Hour),
	}

	userService := &UserService{
		Users: map[string]models.User{
			"dalaoqi": {
				Name: "dalaoqi",
				Folders: map[string]models.Folder{
					"myfolder": {
						Name:  "myfolder",
						Files: map[string]models.File{"myfile1": file1, "myfile2": file2, "myfile3": file3},
					},
				},
			},
		},
	}

	folderService := &FolderService{
		UserService: userService,
	}

	fileService := &FileService{
		UserService:   userService,
		FolderService: folderService,
	}

	testCases := []struct {
		name           string
		userName       string
		folderName     string
		sortFlag       string
		sortOrderFlag  string
		expectedResult []models.File
		expectedError  string
	}{
		{
			name:           "Sort by name in ascending order",
			userName:       "dalaoqi",
			folderName:     "myfolder",
			sortFlag:       "--sort-name",
			sortOrderFlag:  "asc",
			expectedResult: []models.File{file1, file2, file3},
			expectedError:  "",
		},
		{
			name:           "Sort by name in descending order",
			userName:       "dalaoqi",
			folderName:     "myfolder",
			sortFlag:       "--sort-name",
			sortOrderFlag:  "desc",
			expectedResult: []models.File{file3, file2, file1},
			expectedError:  "",
		},
		{
			name:           "Sort by created at in ascending order",
			userName:       "dalaoqi",
			folderName:     "myfolder",
			sortFlag:       "--sort-created",
			sortOrderFlag:  "asc",
			expectedResult: []models.File{file3, file2, file1},
			expectedError:  "",
		},
		{
			name:           "Sort by created at in descending order",
			userName:       "dalaoqi",
			folderName:     "myfolder",
			sortFlag:       "--sort-created",
			sortOrderFlag:  "desc",
			expectedResult: []models.File{file1, file2, file3},
			expectedError:  "",
		},
		{
			name:           "Sort by an invalid flag",
			userName:       "dalaoqi",
			folderName:     "myfolder",
			sortFlag:       "--sort-invalid",
			sortOrderFlag:  "asc",
			expectedResult: []models.File{},
			expectedError:  "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]",
		},
		{
			name:           "Sort by name in invalid order",
			userName:       "dalaoqi",
			folderName:     "myfolder",
			sortFlag:       "--sort-name",
			sortOrderFlag:  "invalid",
			expectedResult: []models.File{},
			expectedError:  "Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]",
		},
		{
			name:           "User doesn't exist",
			userName:       "nonexistentuser",
			folderName:     "myfolder",
			sortFlag:       "--sort-name",
			sortOrderFlag:  "asc",
			expectedResult: []models.File{},
			expectedError:  "Error: The nonexistentuser doesn't exist.",
		},
		{
			name:           "Folder doesn't exist",
			userName:       "dalaoqi",
			folderName:     "nonexistentfolder",
			sortFlag:       "--sort-name",
			sortOrderFlag:  "asc",
			expectedResult: []models.File{},
			expectedError:  "Error: The nonexistentfolder doesn't exist.",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			gotResult, err := fileService.GetFiles(test.userName, test.folderName, test.sortFlag, test.sortOrderFlag)
			if err != nil && err.Error() != test.expectedError {
				t.Errorf("Unexpected error: %s", err)
			}

			if !reflect.DeepEqual(gotResult, test.expectedResult) {
				t.Errorf("Result mismatch, Got: %s, Want: %s", gotResult, test.expectedResult)
			}
		})
	}
}

func TestFileService_DeleteFile(t *testing.T) {
	// Create a test file
	file := models.File{
		Name:      "myfile",
		CreatedAt: time.Now().Add(-time.Hour),
	}

	folder := models.Folder{
		Name:        "myfolder",
		Description: "myfolder",
		CreatedAt:   time.Now().Add(-time.Hour),
		Files: map[string]models.File{
			"myfile": file,
		},
	}

	userService := &UserService{
		Users: map[string]models.User{
			"dalaoqi": {
				Name:    "dalaoqi",
				Folders: map[string]models.Folder{"myfolder": folder},
			},
		},
	}

	folderService := NewFolderService(userService)

	fileService := &FileService{
		UserService:   userService,
		FolderService: folderService,
	}

	testCases := []struct {
		name          string
		userName      string
		folderName    string
		fileName      string
		expectedError string
	}{
		{
			name:          "Delete an existing file",
			userName:      "dalaoqi",
			folderName:    "myfolder",
			fileName:      "myfile",
			expectedError: "",
		},
		{
			name:          "Delete a non-existing file",
			userName:      "dalaoqi",
			folderName:    "myfolder",
			fileName:      "nonexistentfile",
			expectedError: "Error: The nonexistentfile doesn't exist.",
		},
		{
			name:          "Delete a file in a non-existing folder",
			userName:      "dalaoqi",
			folderName:    "nonexistentfolder",
			fileName:      "myfile",
			expectedError: "Error: The nonexistentfolder doesn't exist.",
		},
		{
			name:          "Delete a file in a non-existing user",
			userName:      "nonexistentuser",
			folderName:    "myfolder",
			fileName:      "myfile",
			expectedError: "Error: The nonexistentuser doesn't exist.",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := fileService.DeleteFile(test.userName, test.folderName, test.fileName)

			if err != nil && err.Error() != test.expectedError {
				t.Errorf("Unexpected error: %s", err)
			}

			// Check if the file has been deleted from the folder's files map
			if exists := fileService.Exist(test.userName, test.folderName, test.fileName); exists {
				t.Errorf("File %s still exists in Folder %s for User %s", test.fileName, test.folderName, test.userName)
			}
		})
	}
}
