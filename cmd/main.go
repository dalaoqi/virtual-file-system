package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	vfs "virtual-file-system/internal"
)

var (
	dispatcher *vfs.Dispatcher
)

func init() {
	dispatcher = vfs.NewDispatcher()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("# ")
	// Read input from stdin
	for scanner.Scan() {
		line := scanner.Text()

		// Split the input into individual arguments
		args := strings.Split(line, " ")

		// Execute the appropriate command using the dispatcher
		err := dispatcher.Exec(args)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Print("# ")
	}
}
