package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"

	//
	_ "github.com/lib/pq"
)

// CreateForum create forum
func (db *DataBase) CreateForum(forum *models.Forum) (returnForum models.Forum, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if returnForum, err = db.forumConfirmUnique(tx, forum); err != nil {
		return
	}

	if forum.User, err = db.userCheckID(tx, forum.User); err != nil {
		return
	}

	if returnForum, err = db.createForum(tx, forum); err != nil {
		return
	}

	if err = db.statusAddForum(tx, 1); err != nil {
		return
	}

	err = tx.Commit()
	return
}

// GetForum get forum
func (db *DataBase) GetForum(slug string) (returnForum models.Forum, err error) {
	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if returnForum, err = db.findForumBySlug(tx, slug); err != nil {
		err = re.ErrorUserNotExist()
		return
	}

	err = tx.Commit()
	return
}
