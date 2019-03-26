package database

import (
	"database/sql"
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

// GetPostfull get full post info
func (db *DataBase) GetPostfull(existRelated bool, related string,
	id int) (returnPost models.Postfull, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if returnPost, err = db.postfullGet(tx, existRelated, related, id); err != nil {
		return
	}

	err = tx.Commit()
	return
}
