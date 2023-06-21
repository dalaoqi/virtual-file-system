package models

import "time"

type Folder struct {
	Name        string
	Description string
	CreatedAt   time.Time
}
