package vfs

import (
	"fmt"
	"virtual-file-system/internal/services"
)

type Dispatcher struct {
	userService *services.UserService
}

// NewDispatcher creates a new instance of Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		userService: services.NewUserService(),
	}
}

// Exec executes the command based on the arguments
func (d *Dispatcher) Exec(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Error: Invalid command format")
	}

	switch args[0] {
	case "register":
		// Register a new user using the user service
		err := d.userService.Register(args[1])
		if err != nil {
			return err
		}
		fmt.Printf("Add %v successfully.\n", args[1])
		return nil
	default:
		return fmt.Errorf("Error: Unrecognized command")
	}
}
