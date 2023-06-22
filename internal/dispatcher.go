package vfs

import (
	"fmt"
	"virtual-file-system/internal/services"
)

type Dispatcher struct {
	userService   *services.UserService
	folderService *services.FolderService
	fileService   *services.FileService
}

// NewDispatcher creates a new instance of Dispatcher
func NewDispatcher() *Dispatcher {
	userService := services.NewUserService()
	folderService := services.NewFolderService(userService)
	fileService := services.NewFileService(userService, folderService)
	return &Dispatcher{
		userService:   userService,
		folderService: folderService,
		fileService:   fileService,
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

		err := d.folderService.CreateFolder(userName, folderName, description)
		if err != nil {
			return err
		}
		fmt.Printf("Create %v successfully.\n", folderName)
		return nil
	case "list-folders":
		if len(args) < 2 {
			return fmt.Errorf("Error: Insufficient arguments")
		}
		userName := args[1]
		sortFlag := "--sort-name"
		sortOrderFlag := "asc"

		if len(args) > 2 {
			// Check if the sort flag is provided
			sortFlag = args[2]

			if len(args) > 3 {
				// Check if the sort order is provided
				sortOrderFlag = args[3]

			}
		}

		folders, err := d.folderService.GetFolders(userName, sortFlag, sortOrderFlag)
		if err != nil {
			return err
		}

		for _, folder := range folders {
			createdAt := folder.CreatedAt.Format("2006-01-02 15:04:05")
			fmt.Printf("%s %s %s %s\n", folder.Name, folder.Description, createdAt, userName)
		}

		return nil
	case "delete-folder":
		if len(args) < 3 {
			return fmt.Errorf("Error: Insufficient arguments")
		}
		userName := args[1]
		folderName := args[2]

		err := d.folderService.DeleteFolder(userName, folderName)
		if err != nil {
			return err
		}
		fmt.Printf("Delete %v successfully.\n", folderName)
		return nil
	case "rename-folder":
		if len(args) < 4 {
			return fmt.Errorf("Error: Insufficient arguments")
		}
		userName := args[1]
		folderName := args[2]
		newFolderName := args[3]

		err := d.folderService.RenameFolder(userName, folderName, newFolderName)
		if err != nil {
			return err
		}
		fmt.Printf("Rename %s to %s successfully.\n", folderName, newFolderName)
		return nil
	case "create-file":
		if len(args) < 4 {
			return fmt.Errorf("Error: Insufficient arguments")
		}
		userName := args[1]
		folderName := args[2]
		fileName := args[3]
		description := ""
		if len(args) > 4 {
			description = args[4]
		}

		err := d.fileService.CreateFile(userName, folderName, fileName, description)
		if err != nil {
			return err
		}
		fmt.Printf("Create %v in %v/%v successfully.\n", fileName, userName, folderName)
		return nil
	case "list-files":
		if len(args) < 3 {
			return fmt.Errorf("Error: Insufficient arguments")
		}
		userName := args[1]
		folderName := args[2]
		sortFlag := "--sort-name"
		sortOrderFlag := "asc"

		if len(args) > 3 {
			// Check if the sort flag is provided
			sortFlag = args[3]

			if len(args) > 4 {
				// Check if the sort order is provided
				sortOrderFlag = args[4]
			}
		}

		files, err := d.fileService.GetFiles(userName, folderName, sortFlag, sortOrderFlag)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			fmt.Println("Warning: The folder is empty.")
			return nil
		}

		for _, file := range files {
			createdAt := file.CreatedAt.Format("2006-01-02 15:04:05")
			fmt.Printf("%s %s %s %s %s\n", file.Name, file.Description, createdAt, folderName, userName)
		}
		return nil
	default:
		return fmt.Errorf("Error: Unrecognized command")
	}
}
