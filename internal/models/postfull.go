package models

// Postfull model
type Postfull struct {
	Author *User   `json:"author,omitempty" db:"author"`
	Forum  *Forum  `json:"forum,omitempty" db:"forum"`
	Post   *Post   `json:"post,omitempty" db:"post"`
	Thread *Thread `json:"thread,omitempty" db:"thread"`
}
