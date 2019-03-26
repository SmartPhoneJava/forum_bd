package models

import (
	"fmt"
	"time"
)

// Post model
type Post struct {
	Author   string    `json:"author" db:"author"`
	Created  time.Time `json:"created" db:"created"`
	Forum    string    `json:"forum" db:"forum"`
	ID       int       `json:"id" db:"id"`
	IsEdited bool      `json:"isEdited" db:"isEdited"`
	Message  string    `json:"message" db:"message"`
	Parent   int       `json:"parent" db:"parent"`
	Thread   int       `json:"thread" db:"thread"`
	Path     string    `json:"path" db:"path"`
}

// Print for debug
func (post *Post) Print() {
	fmt.Println("-------Post-------")
	fmt.Println("--ID:", post.ID)
	fmt.Println("--Parent:", post.Parent)
	fmt.Println("--Path:", post.Path)
	fmt.Println("--Created:", post.Created)
	fmt.Println("--IsEdited:", post.IsEdited)
	fmt.Println("--Thread:", post.Thread)
}
