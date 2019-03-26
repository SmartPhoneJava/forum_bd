package models

// Result model
type Result struct {
	Place   string `json:"-"`
	Success bool   `json:"-"`
	Message string `json:"message"`
}
