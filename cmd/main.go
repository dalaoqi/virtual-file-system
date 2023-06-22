package main

import (
	"bufio"
	"fmt"
	"os"
	"virtual-file-system/internal/services"
	"virtual-file-system/internal/utils"
)

var (
	dispatcher *services.Dispatcher
)

func init() {
	dispatcher = services.NewDispatcher()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("# ")
	// Read input from stdin
	for scanner.Scan() {
		line := scanner.Text()

		// Split the input into individual arguments
		args := utils.SplitArguments(line)
		// Execute the appropriate command using the dispatcher
		err := dispatcher.Exec(args)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Print("# ")
	}
}
