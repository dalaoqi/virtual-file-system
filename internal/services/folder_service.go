package services

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"virtual-file-system/internal/models"
)

type FolderService struct {
	Folders map[string]map[string]models.Folder
}

// NewFolderService creates a new instance of FolderService
func NewFolderService() *FolderService {
	return &FolderService{
		Folders: make(map[string]map[string]models.Folder),
	}
}

func (s *FolderService) CreateFolder(userName, folderName, description string) error {
	userName = strings.ToLower(userName)
	folderName = strings.ToLower(folderName)

	// Check if the user's folder map exists
	if s.Folders[userName] == nil {
		s.Folders[userName] = make(map[string]models.Folder)
	} else {
		// Check if the folder name already exists for the user
		for _, folder := range s.Folders[userName] {
			if folder.Name == folderName {
				return fmt.Errorf("Error: The folder name '%s' already exists for the user '%s'.", folderName, userName)
			}
		}
	}

	// Create the new folder
	s.Folders[userName][folderName] = models.Folder{
		Name:        folderName,
		Owner:       userName,
		Description: description,
		CreatedAt:   time.Now(),
	}

	return nil
}

func (s *FolderService) GetFolders(userName, sortBy, sortOrder string) ([]models.Folder, error) {
	userName = strings.ToLower(userName)

	folders, ok := s.Folders[userName]
	if !ok {
		return []models.Folder{}, nil
	}

	// Convert map to slice for sorting
	folderList := make([]models.Folder, 0)
	for _, folder := range folders {
		folderList = append(folderList, folder)
	}

	// Sort the folders based on the provided flags
	switch sortBy {
	case "name":
		if sortOrder == "asc" {
			sort.SliceStable(folderList, func(i, j int) bool {
				return folderList[i].Name < folderList[j].Name
			})
		} else {
			sort.SliceStable(folderList, func(i, j int) bool {
				return folderList[i].Name > folderList[j].Name
			})
		}
	case "created":
		if sortOrder == "asc" {
			sort.SliceStable(folderList, func(i, j int) bool {
				return folderList[i].CreatedAt.Before(folderList[j].CreatedAt)
			})
		} else {
			sort.SliceStable(folderList, func(i, j int) bool {
				return folderList[i].CreatedAt.After(folderList[j].CreatedAt)
			})
		}
	}
	return folderList, nil
}

func (s *FolderService) DeleteFolder(userName, folderName string) {
	userName = strings.ToLower(userName)
	folderName = strings.ToLower(folderName)

	delete(s.Folders[userName], folderName)
	return
}

func (s *FolderService) Exist(userName, folderName string) bool {
	userName = strings.ToLower(userName)
	folderName = strings.ToLower(folderName)

	folders, ok := s.Folders[userName]
	if !ok {
		return false
	}
	_, exist := folders[folderName]
	return exist
}
