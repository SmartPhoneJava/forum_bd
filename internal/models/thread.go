package models

// Thread model
type Thread struct {
	Author  string `json:"author,omitempty" db:"author"`
	Created string `json:"created,omitempty" db:"created"`
	Forum   string `json:"forum,omitempty" db:"forum"`
	ID      int    `json:"id,omitempty" db:"id"`
	Message string `json:"message,omitempty" db:"message"`
	Slug    string `json:"slug,omitempty" db:"slug"`
	Title   string `json:"title,omitempty" db:"title"`
	Votes   int    `json:"votes,omitempty" db:"votes"`
}
