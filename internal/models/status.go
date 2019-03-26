package models

// Status model
type Status struct {
	Forum  int `json:"forum,omitempty" db:"forum"`
	Post   int `json:"post,omitempty" db:"post"`
	Thread int `json:"thread,omitempty" db:"thread"`
	User   int `json:"user,omitempty" db:"user"`
}
