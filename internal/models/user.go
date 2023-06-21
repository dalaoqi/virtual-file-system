package models

// User represents a user in the system
type User struct {
	Name    string
	Folders map[string]Folder
}
