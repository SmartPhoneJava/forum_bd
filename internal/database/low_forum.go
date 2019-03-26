package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"strconv"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) createForum(tx *sql.Tx, forum *models.Forum) (createdForum models.Forum, err error) {

	query := `INSERT INTO Forum(slug, title, user_nickname ) VALUES
						 	($1, $2, $3) 
						 `
	queryAddForumReturning(&query)
	row := tx.QueryRow(query, forum.Slug, forum.Title, forum.User)
	createdForum, err = forumScan(row)
	return
}

func (db *DataBase) forumCheckID(tx *sql.Tx, oldNickname string) (newNickname string, err error) {
	var thatForum models.Forum
	if thatForum, err = db.findForumBySlug(tx, oldNickname); err != nil {
		err = re.ErrorForumNotExist()
		return
	}
	newNickname = thatForum.Slug
	return
}

func (db DataBase) forumConfirmUnique(tx *sql.Tx, forum *models.Forum) (foundForum models.Forum, err error) {

	if foundForum, err = db.findForumBySlug(tx, forum.Slug); err != nil && err != sql.ErrNoRows {
		return
	}
	err = nil

	if foundForum.Slug != "" {
		err = re.ErrorForumSlugIsTaken()
		return
	}

	return
}

func findForum(tx *sql.Tx, queryAdd string, arg string) (foundForum models.Forum, err error) {

	query := querySelectForum() + queryAdd

	row := tx.QueryRow(query, arg)
	foundForum, err = forumScan(row)
	return
}

// findForumBySlug
func (db DataBase) findForumBySlug(tx *sql.Tx, taken string) (foundForum models.Forum, err error) {

	query := `where lower(slug) like lower($1)`
	foundForum, err = findForum(tx, query, taken)
	return
}

// forumUpdateThreads
func (db DataBase) forumUpdateThreads(tx *sql.Tx, slug string) (err error) {
	query := `UPDATE Forum set threads=threads+1 where slug=$1`
	_, err = tx.Exec(query, slug)
	return
}

// forumUpdatePosts
func (db DataBase) forumUpdatePosts(tx *sql.Tx, slug string, amount int) (err error) {
	query := `UPDATE Forum set posts=posts+` + strconv.Itoa(amount) + ` where slug=$1`
	_, err = tx.Exec(query, slug)
	return
}

// query add returning
func queryAddForumReturning(query *string) {
	*query += ` returning slug, title, user_nickname, posts, threads `
}

// queryAddThreadReturning add returning for insert,update etc
func querySelectForum() string {
	return ` SELECT F.slug, F.title, F.user_nickname, F.posts, F.threads FROM Forum as F `
}

// scan to model Forum
func forumScan(row *sql.Row) (foundForum models.Forum, err error) {
	foundForum = models.Forum{}
	err = row.Scan(&foundForum.Slug, &foundForum.Title,
		&foundForum.User, &foundForum.Posts, &foundForum.Threads)
	foundForum.Print()
	return
}
