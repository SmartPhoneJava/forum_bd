package models

// User model
type User struct {
	ID       int    `json:"-" db:"id"`
	About    string `json:"about,omitempty" db:"about"`
	Email    string `json:"email,omitempty" db:"email"`
	Fullname string `json:"fullname,omitempty" db:"fullname"`
	Nickname string `json:"nickname,omitempty" db:"nickname"`
}

// Print for debug
func (user *User) Print() {
	/*
		fmt.Println("-------User-------")
		fmt.Println("--About:", user.About)
		fmt.Println("--Email:", user.Email)
		fmt.Println("--Fullname:", user.Fullname)
		fmt.Println("--Nickname:", user.Nickname)
	*/
}
