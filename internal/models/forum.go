package models

import "fmt"

// Forum model
type Forum struct {
	Posts   int    `json:"posts" db:"posts"`
	Threads int    `json:"threads" db:"threads"`
	Slug    string `json:"slug" db:"slug"`
	Title   string `json:"title" db:"title"`
	User    string `json:"user" db:"user_nickname"`
}

// Print for debug
func (forum *Forum) Print() {
	fmt.Println("-------Forum-------")
	fmt.Println("--Posts:", forum.Posts)
	fmt.Println("--Threads:", forum.Threads)
	fmt.Println("--Slug:", forum.Slug)
	fmt.Println("--Title:", forum.Title)
	fmt.Println("--User:", forum.User)
}

/*
 CREATE Table Forum (
        posts int default 0,
        slug SERIAL PRIMARY KEY,
        threads int,
        title varchar(60) not null,
        user_nickname varchar(80) not null
    );
*/
