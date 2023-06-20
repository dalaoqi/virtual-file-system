package models

import "time"

type Folder struct {
	Name        string
	Owner       string
	Description string
	CreatedAt   time.Time
}
