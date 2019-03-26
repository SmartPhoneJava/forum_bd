package database

import (
	"database/sql"
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) statusGet(tx *sql.Tx) (status models.Status, err error) {

	query := `select forum, post, thread, users from Status `
	status = models.Status{}
	err = tx.QueryRow(query).Scan(&status.Forum, &status.Post,
		&status.Thread, &status.User)

	return
}

func (db *DataBase) clearDataBase(tx *sql.Tx) (status models.Status, err error) {

	query := deleteTables()
	err = tx.QueryRow(query).Scan(&status.Forum, &status.Post,
		&status.Thread, &status.User)

	return
}

func (db *DataBase) statusAddForum(tx *sql.Tx, add int) (err error) {

	query := `UPDATE Status set forum = forum + $1`
	_, err = tx.Exec(query, add)
	return
}

func (db *DataBase) statusAddPost(tx *sql.Tx, add int) (err error) {

	query := `UPDATE Status set post = post + $1`
	_, err = tx.Exec(query, add)
	return
}

func (db *DataBase) statusAddThread(tx *sql.Tx, add int) (err error) {

	query := `UPDATE Status set thread = thread + $1`
	_, err = tx.Exec(query, add)
	return
}

func (db *DataBase) statusAddUser(tx *sql.Tx, add int) (err error) {

	query := `UPDATE Status set users = users + $1`
	_, err = tx.Exec(query, add)
	return
}

func deleteTables() string {
	return `
	DELETE FROM Vote;
	DELETE FROM Post;
	DELETE FROM Thread;
	DELETE FROM Forum;
	DELETE FROM UserForum;
	DELETE FROM Status;
	INSERT INTO Status(Post) VALUES (0) 
	RETURNING forum, post, thread, users
    `
}
