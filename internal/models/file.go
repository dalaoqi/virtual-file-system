package models

import "time"

// File represents a file in the system
type File struct {
	Name        string
	Description string
	CreatedAt   time.Time
}
