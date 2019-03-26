package models

import "fmt"

// Vote model
type Vote struct {
	Author   string `json:"nickname" db:"author"`
	Voice    int    `json:"voice" db:"voice"`
	Thread   int    `json:"-" db:"thread"`
	IsEdited bool   `json:"-" db:"isEdited"`
}

// Print for debug
func (vote *Vote) Print() {
	fmt.Println("-------Vote-------")
	fmt.Println("--Author					:", vote.Author)
	fmt.Println("--Voice					:", vote.Voice)
	fmt.Println("--Thread					:", vote.Thread)
	fmt.Println("--IsEdited				:", vote.IsEdited)
}
