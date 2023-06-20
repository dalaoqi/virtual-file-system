package vfs

import (
	"fmt"
	"virtual-file-system/internal/services"
	"virtual-file-system/internal/utils"
)

type Dispatcher struct {
	userService   *services.UserService
	folderService *services.FolderService
}

// NewDispatcher creates a new instance of Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		userService:   services.NewUserService(),
		folderService: services.NewFolderService(),
	}
}

// Exec executes the command based on the arguments
func (d *Dispatcher) Exec(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Error: Insufficient arguments")
	}

	switch args[0] {
	case "register":
		userName := args[1]

		// Check if the user already exists
		if d.userService.Exist(userName) {
			return fmt.Errorf("Error: The %v has already existed.", userName)
		}

		// Check if the name contains invalid characters
		if utils.ExistInvalidChars(userName) {
			return fmt.Errorf("Error: The %v contains invalid chars.", userName)
		}

		// Register a new user using the user service
		err := d.userService.Register(userName)
		if err != nil {
			return err
		}
		fmt.Printf("Add %v successfully.\n", userName)
		return nil
	case "create-folder":
		if len(args) < 3 {
			return fmt.Errorf("Error: Insufficient arguments")
		}
		userName := args[1]
		folderName := args[2]
		description := ""
		if len(args) > 3 {
			description = args[3]
		}

		// Check if the user already exists
		if !d.userService.Exist(userName) {
			return fmt.Errorf("Error: The %v doesn't exist.", userName)
		}

		// Check if the folder name contains invalid characters
		if utils.ExistInvalidChars(folderName) {
			return fmt.Errorf("Error: The %v contains invalid chars.", folderName)
		}

		// Check if the folder name already exists for the user
		if d.folderService.Exist(userName, folderName) {
			return fmt.Errorf("Error: The folder name '%s' already exists for the user '%s'.", folderName, userName)
		}
		err := d.folderService.CreateFolder(userName, folderName, description)
		if err != nil {
			return err
		}
		fmt.Printf("Create %v successfully.\n", folderName)
		return nil
	default:
		return fmt.Errorf("Error: Unrecognized command")
	}
}
