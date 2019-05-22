package models

// Post model
type Post struct {
	Author   string `json:"author,omitempty" db:"author"`
	Created  string `json:"created,omitempty" db:"created"`
	Forum    string `json:"forum,omitempty" db:"forum"`
	ID       int    `json:"id,omitempty" db:"id"`
	IsEdited bool   `json:"isEdited,omitempty" db:"isEdited"`
	Message  string `json:"message,omitempty" db:"message"`
	Parent   int    `json:"parent,omitempty" db:"parent"`
	Thread   int    `json:"thread,omitempty" db:"thread"`
	Path     string `json:"path,omitempty" db:"path"`
}

// Print for debug
func (post *Post) Print() {
	/*
		fmt.Println("-------Post-------")
		fmt.Println("--ID:", post.ID)
		fmt.Println("--Parent:", post.Parent)
		fmt.Println("--Path:", post.Path)
		fmt.Println("--Created:", post.Created)
		fmt.Println("--IsEdited:", post.IsEdited)
		fmt.Println("--Thread:", post.Thread)
	*/
}
