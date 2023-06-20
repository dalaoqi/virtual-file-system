package services

import (
	"fmt"
	"strings"
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
	}

	return nil
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
