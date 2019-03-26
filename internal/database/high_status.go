package database

import (
	"database/sql"
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

// GetStatus get all info about database
func (db *DataBase) GetStatus() (returnStatus models.Status, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if returnStatus, err = db.statusGet(tx); err != nil {
		return
	}

	err = tx.Commit()
	return
}

// Clear handle deleting all info
func (db *DataBase) Clear() (status models.Status, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if status, err = db.clearDataBase(tx); err != nil {
		return
	}

	err = tx.Commit()
	return
}
