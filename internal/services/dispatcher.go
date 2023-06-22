package services

import (
	"fmt"
)

type Dispatcher struct {
	userService   *UserService
	folderService *FolderService
	fileService   *FileService
}

// NewDispatcher creates a new instance of Dispatcher
func NewDispatcher() *Dispatcher {
	userService := NewUserService()
	folderService := NewFolderService(userService)
	fileService := NewFileService(userService, folderService)
	return &Dispatcher{
		userService:   userService,
		folderService: folderService,
		fileService:   fileService,
	}
}

// Exec executes the command based on the arguments
func (d *Dispatcher) Exec(args []string) error {
	switch args[0] {
	case "register":
		if len(args) < 2 {
			return fmt.Errorf("Error: Insufficient arguments\nUsage: register [username]")
		}
		userName := args[1]

		// Register a new user using the user service
		err := d.userService.Register(userName)
		if err != nil {
			return err
		}

		fmt.Printf("Add %s successfully.\n", userName)
		return nil
	case "create-folder":
		if len(args) < 3 {
			return fmt.Errorf("Error: Insufficient arguments\nUsage: create-folder [username] [foldername] [description]?")
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
		fmt.Printf("Create %s successfully.\n", folderName)
		return nil
	case "list-folders":
		if len(args) < 2 {
			return fmt.Errorf("Error: Insufficient arguments\nUsage: list-folders [username] [--sort-name|--sort-created] [asc|desc]")
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
			return fmt.Errorf("Error: Insufficient arguments\nUsage: delete-folder [username] [foldername]")
		}
		userName := args[1]
		folderName := args[2]

		err := d.folderService.DeleteFolder(userName, folderName)
		if err != nil {
			return err
		}
		fmt.Printf("Delete %s successfully.\n", folderName)
		return nil
	case "rename-folder":
		if len(args) < 4 {
			return fmt.Errorf("Error: Insufficient arguments\nUsage: rename-folder [username] [foldername] [new-folder-name]")
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
			return fmt.Errorf("Error: Insufficient arguments\nUsage: create-file [username] [foldername] [filename] [description]?")
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
		fmt.Printf("Create %s in %s/%s successfully.\n", fileName, userName, folderName)
		return nil
	case "list-files":
		if len(args) < 3 {
			return fmt.Errorf("Error: Insufficient arguments\nUsage: list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
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
	case "delete-file":
		if len(args) < 4 {
			return fmt.Errorf("Error: Insufficient arguments\nUsage: delete-file [username] [foldername] [filename]")
		}

		userName := args[1]
		folderName := args[2]
		fileName := args[3]

		err := d.fileService.DeleteFile(userName, folderName, fileName)
		if err != nil {
			return err
		}

		fmt.Printf("Delete %s in %s/%s successfully.\n", fileName, userName, folderName)
		return nil
	default:
		return fmt.Errorf("Error: Unrecognized command")
	}
}
