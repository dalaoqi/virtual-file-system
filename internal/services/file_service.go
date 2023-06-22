package services

import (
	"fmt"
	"sort"
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

func (s *FileService) GetFiles(userName, folderName, sortFlag, sortOrderFlag string) ([]models.File, error) {
	lowerUserName := strings.ToLower(userName)
	lowerFolderName := strings.ToLower(folderName)

	// Check if the user exists
	if !s.UserService.Exist(lowerUserName) {
		return []models.File{}, fmt.Errorf("Error: The %v doesn't exist.", userName)
	}

	// Check if the folder exists for the user
	if !s.FolderService.Exist(lowerUserName, lowerFolderName) {
		return []models.File{}, fmt.Errorf("Error: The %s doesn't exist.", folderName)
	}

	// Convert map to slice for sorting
	fileList := make([]models.File, 0)
	for _, file := range s.UserService.Users[lowerUserName].Folders[lowerFolderName].Files {
		fileList = append(fileList, file)
	}

	if len(fileList) == 0 {
		return fileList, nil
	}

	// Sort the files based on the provided flags
	switch sortFlag {
	case "--sort-name":
		if sortOrderFlag == "asc" {
			sort.SliceStable(fileList, func(i, j int) bool {
				return fileList[i].Name < fileList[j].Name
			})
		} else if sortOrderFlag == "desc" {
			sort.SliceStable(fileList, func(i, j int) bool {
				return fileList[i].Name > fileList[j].Name
			})
		} else {
			return []models.File{}, fmt.Errorf("Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
		}
	case "--sort-created":
		if sortOrderFlag == "asc" {
			sort.SliceStable(fileList, func(i, j int) bool {
				return fileList[i].CreatedAt.Before(fileList[j].CreatedAt)
			})
		} else if sortOrderFlag == "desc" {
			sort.SliceStable(fileList, func(i, j int) bool {
				return fileList[i].CreatedAt.After(fileList[j].CreatedAt)
			})
		} else {
			return []models.File{}, fmt.Errorf("Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
		}
	default:
		return []models.File{}, fmt.Errorf("Usage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
	}

	return fileList, nil
}

func (s *FileService) DeleteFile(userName, folderName, fileName string) error {
	lowerUserName := strings.ToLower(userName)
	lowerFolderName := strings.ToLower(folderName)
	lowerFileName := strings.ToLower(fileName)

	// Check if the user exists
	if !s.UserService.Exist(lowerUserName) {
		return fmt.Errorf("Error: The %s doesn't exist.", userName)
	}

	// Check if the folder exists for the user
	if !s.FolderService.Exist(lowerUserName, lowerFolderName) {
		return fmt.Errorf("Error: The %s doesn't exist.", folderName)
	}

	// Check if the file exists in the folder
	if !s.Exist(lowerUserName, lowerFolderName, lowerFileName) {
		return fmt.Errorf("Error: The %s doesn't exist.", fileName)
	}

	// Delete the file from the folder
	delete(s.UserService.Users[lowerUserName].Folders[lowerFolderName].Files, lowerFileName)

	return nil
}

func (s *FileService) Exist(userName, folderName, fileName string) bool {
	file, exist := s.UserService.Users[userName].Folders[folderName].Files[fileName]
	if !exist {
		return false
	}
	return file.Name == fileName
}
