package database

import (
	"database/sql"
	"time"

	"github.com/SmartPhoneJava/forum_bd/internal/models"
	re "github.com/SmartPhoneJava/forum_bd/internal/return_errors"

	//
	_ "github.com/lib/pq"
)

// CreateThread handle thread creation
func (db *DataBase) CreateThread(thread *models.Thread) (returnThread models.Thread, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if returnThread, err = db.threadConfirmUnique(tx, thread); err != nil {
		err = re.ErrorThreadConflict()
		return
	}

	// if thread.Author, err = db.userCheckID(tx, thread.Author); err != nil {
	// 	return
	// }

	// debug("forumCheckID:", thread.Forum)
	// if thread.Forum, err = db.forumCheckID(tx, thread.Forum); err != nil {
	// 	return
	// }
	debug("forumCheckID1:", thread.Forum)
	if returnThread, err = db.threadCreate(tx, thread); err != nil {
		return
	}

	// if err = db.forumUpdateThreads(tx, thread.Forum); err != nil {
	// 	return
	// }

	// if err = db.statusAddThread(tx, 1); err != nil {
	// 	return
	// }

	if err = tx.Commit(); err != nil {
		return
	}

	var tx1 *sql.Tx
	if tx1, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx1.Rollback()

	if erro := db.userInForumCreate(thread.Author, thread.Forum); erro != nil {
		debug("erro:", erro.Error())
		return
	}
	tx1.Commit()

	return
}

// UpdateThread handle thread update
func (db *DataBase) UpdateThread(thread *models.Thread,
	slug string) (returnThread models.Thread, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	// if returnThread, err = db.threadConfirmUnique(tx, thread); err != nil {
	// 	return
	// }

	if returnThread, err = db.threadUpdate(tx, thread, slug); err != nil {
		return
	}
	err = tx.Commit()
	return
}

// GetThreads get threads
func (db *DataBase) GetThreads(slug string, limit int, existLimit bool, t time.Time, existTime bool, desc bool) (returnThreads []models.Thread, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if _, err = db.findForumBySlug(tx, slug); err != nil {
		err = re.ErrorForumNotExist()
		return
	}

	if returnThreads, err = db.threadsGet(tx, slug, limit, existLimit, t, existTime, desc); err != nil {
		return
	}

	err = tx.Commit()
	return
}

// GetThread get thread
func (db *DataBase) GetThread(slug string) (returnThread models.Thread, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if returnThread, err = db.threadFindByIDorSlug(tx, slug); err != nil {
		return
	}

	err = tx.Commit()
	return
}
