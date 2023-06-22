package services

import (
	"fmt"
	"strings"
	"time"
	"virtual-file-system/internal/models"
	"virtual-file-system/internal/utils"
)

type FileService struct {
	UserService   *UserService
	FolderService *FolderService
}

// NewFileService creates a new instance of FileService
func NewFileService(userService *UserService, folderService *FolderService) *FileService {
	return &FileService{
		UserService:   userService,
		FolderService: folderService,
	}
}

func (s *FileService) CreateFile(userName, folderName, fileName, description string) error {
	lowerUserName := strings.ToLower(userName)
	lowerFolderName := strings.ToLower(folderName)
	lowerFileName := strings.ToLower(fileName)

	// Check if the user exists
	if !s.UserService.Exist(lowerUserName) {
		return fmt.Errorf("Error: The %v doesn't exist.", userName)
	}

	// Check if the folder exists for the user
	if !s.FolderService.Exist(lowerUserName, lowerFolderName) {
		return fmt.Errorf("Error: The %s doesn't exist.", folderName)
	}

	// Check if the file name contains invalid characters
	if utils.ExistInvalidChars(lowerFileName) {
		return fmt.Errorf("Error: The %v contains invalid chars.", fileName)
	}

	// Check if the file name already exists in the folder
	if s.Exist(lowerUserName, lowerFolderName, lowerFileName) {
		return fmt.Errorf("Error: The %s has already existed in the %s.", fileName, folderName)
	}

	// Create the new file
	folder := s.UserService.Users[lowerUserName].Folders[lowerFolderName]

	if folder.Files == nil {
		folder.Files = make(map[string]models.File)
	}
	folder.Files[lowerFileName] = models.File{
		Name:        lowerFileName,
		Description: description,
		CreatedAt:   time.Now(),
	}
	s.UserService.Users[lowerUserName].Folders[lowerFolderName] = folder
	return nil
}

func (s *FileService) Exist(userName, folderName, fileName string) bool {
	file, exist := s.UserService.Users[userName].Folders[folderName].Files[fileName]
	if !exist {
		return false
	}
	return file.Name == fileName
}
