package services

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"virtual-file-system/internal/models"
	"virtual-file-system/internal/utils"
)

type FolderService struct {
	UserService *UserService
}

// NewFolderService creates a new instance of FolderService
func NewFolderService(userService *UserService) *FolderService {
	return &FolderService{
		UserService: userService,
	}
}

func (s *FolderService) CreateFolder(userName, folderName, description string) error {
	lowerUserName := strings.ToLower(userName)
	lowerFolderName := strings.ToLower(folderName)

	// Check if the user already exists
	if !s.UserService.Exist(lowerUserName) {
		return fmt.Errorf("Error: The %v doesn't exist.", userName)
	}

	// Check if the folder name contains invalid characters
	if utils.ExistInvalidChars(lowerFolderName) {
		return fmt.Errorf("Error: The %v contains invalid chars.", folderName)
	}

	// Check if the folder name already exists for the user
	if s.Exist(lowerUserName, lowerFolderName) {
		return fmt.Errorf("Error: The %s has already existed.", folderName)
	}

	// Create the new folder
	user := s.UserService.Users[lowerUserName]
	if user.Folders == nil {
		user.Folders = make(map[string]models.Folder)
	}
	user.Folders[lowerFolderName] = models.Folder{
		Name:        lowerFolderName,
		Description: description,
		CreatedAt:   time.Now(),
	}

	s.UserService.Users[lowerUserName] = user
	return nil
}

func (s *FolderService) GetFolders(userName, sortFlag, sortOrderFlag string) ([]models.Folder, error) {
	lowerUserName := strings.ToLower(userName)

	// Check if the user exists
	if !s.UserService.Exist(lowerUserName) {
		return []models.Folder{}, fmt.Errorf("Error: The %v doesn't exist.", userName)
	}

	// Convert map to slice for sorting
	folderList := make([]models.Folder, 0)
	for _, folder := range s.UserService.Users[lowerUserName].Folders {
		folderList = append(folderList, folder)
	}

	if len(folderList) == 0 {
		return folderList, fmt.Errorf("Warning: The %s doesn't have any folders.\n", userName)
	}

	// Sort the folders based on the provided flags
	switch sortFlag {
	case "--sort-name":
		if sortOrderFlag == "asc" {
			sort.SliceStable(folderList, func(i, j int) bool {
				return folderList[i].Name < folderList[j].Name
			})
		} else if sortOrderFlag == "desc" {
			sort.SliceStable(folderList, func(i, j int) bool {
				return folderList[i].Name > folderList[j].Name
			})
		} else {
			return []models.Folder{}, fmt.Errorf("Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]")
		}
	case "--sort-created":
		if sortOrderFlag == "asc" {
			sort.SliceStable(folderList, func(i, j int) bool {
				return folderList[i].CreatedAt.Before(folderList[j].CreatedAt)
			})
		} else if sortOrderFlag == "desc" {
			sort.SliceStable(folderList, func(i, j int) bool {
				return folderList[i].CreatedAt.After(folderList[j].CreatedAt)
			})
		} else {
			return []models.Folder{}, fmt.Errorf("Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]")
		}
	default:
		return []models.Folder{}, fmt.Errorf("Usage: list-folders [username] [--sort-name|--sort-created] [asc|desc]")
	}

	return folderList, nil
}

func (s *FolderService) DeleteFolder(userName, folderName string) error {
	lowerUserName := strings.ToLower(userName)
	lowerFolderName := strings.ToLower(folderName)

	// Check if the user exists
	if !s.UserService.Exist(lowerUserName) {
		return fmt.Errorf("Error: The %v doesn't exist.", userName)
	}

	// Check if the folder exists for the user
	if !s.Exist(lowerUserName, lowerFolderName) {
		return fmt.Errorf("Error: The %s doesn't exist", folderName)
	}

	delete(s.UserService.Users[lowerUserName].Folders, lowerFolderName)
	return nil
}

func (s *FolderService) RenameFolder(userName, folderName, newFolderName string) error {
	lowerUserName := strings.ToLower(userName)
	lowerFolderName := strings.ToLower(folderName)
	lowerNewFolderName := strings.ToLower(newFolderName)

	// Check if the user exists
	if !s.UserService.Exist(lowerUserName) {
		return fmt.Errorf("Error: The %v doesn't exist.", userName)
	}

	// Check if the folder exists for the user
	if !s.Exist(lowerUserName, lowerFolderName) {
		return fmt.Errorf("Error: The %v doesn't exist", folderName)
	}

	// Check if the new folder name contains invalid characters
	if utils.ExistInvalidChars(lowerNewFolderName) {
		return fmt.Errorf("Error: The %v contains invalid chars.", newFolderName)
	}

	// Check if the folder exists for the user
	if s.Exist(lowerUserName, lowerNewFolderName) {
		return fmt.Errorf("Error: The %v has already existed.", newFolderName)
	}

	folder := s.UserService.Users[lowerUserName].Folders[lowerFolderName]
	folder.Name = lowerNewFolderName
	s.UserService.Users[lowerUserName].Folders[lowerNewFolderName] = folder
	delete(s.UserService.Users[lowerUserName].Folders, lowerFolderName)
	return nil
}

func (s *FolderService) Exist(userName, folderName string) bool {
	folder, exist := s.UserService.Users[userName].Folders[folderName]
	if !exist {
		return false
	}
	return folder.Name == folderName
}
