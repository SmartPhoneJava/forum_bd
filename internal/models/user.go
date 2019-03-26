package models

import (
	"fmt"
)

// User model
type User struct {
	ID       int    `json:"-" db:"id"`
	About    string `json:"about" db:"about"`
	Email    string `json:"email" db:"email"`
	Fullname string `json:"fullname" db:"fullname"`
	Nickname string `json:"nickname" db:"nickname"`
}

// Print for debug
func (user *User) Print() {
	fmt.Println("-------User-------")
	fmt.Println("--About:", user.About)
	fmt.Println("--Email:", user.Email)
	fmt.Println("--Fullname:", user.Fullname)
	fmt.Println("--Nickname:", user.Nickname)
}
