package database

import (
	"database/sql"
	"escapade/internal/models"
	"fmt"
	"strings"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) postfullGet(tx *sql.Tx, existRelated bool, related string, id int) (fullpost models.Postfull, err error) {

	query := `
	select id, author, created, forum,
	message, thread, parent, path, isEdited
		 from Post 
		 where id = $1
	`
	fullpost.Post = &models.Post{}
	if err = tx.QueryRow(query, id).Scan(&fullpost.Post.ID, &fullpost.Post.Author, &fullpost.Post.Created,
		&fullpost.Post.Forum, &fullpost.Post.Message, &fullpost.Post.Thread, &fullpost.Post.Parent,
		&fullpost.Post.Path, &fullpost.Post.IsEdited); err != nil {
		return
	}

	if strings.Contains(related, "user") {
		query = querySelectUser() + ` join Post on UserForum.nickname=Post.author
			 where Post.id = $1
		`
		fullpost.Author = &models.User{}

		if *fullpost.Author, err = userScan(tx.QueryRow(query, id)); err != nil {
			return
		}
	}
	if strings.Contains(related, "thread") {
		query = querySelectThread() + ` join Post as P on T.id=P.thread
			 where P.id = $1
		`
		fmt.Println("query Thread:" + query)
		fullpost.Thread = &models.Thread{}

		if *fullpost.Thread, err = threadScan(tx.QueryRow(query, id)); err != nil {
			return
		}
	}
	if strings.Contains(related, "forum") {
		query = querySelectForum() + ` join Post as P on lower(F.slug)= lower(P.forum)
			 where P.id = $1
		`
		fullpost.Forum = &models.Forum{}

		if *fullpost.Forum, err = forumScan(tx.QueryRow(query, id)); err != nil {
			return
		}
	}

	return
}
