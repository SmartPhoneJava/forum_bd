package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"fmt"
	"strconv"
	"time"

	//
	_ "github.com/lib/pq"
)

// createThread create thread
func (db *DataBase) threadCreate(tx *sql.Tx, thread *models.Thread) (createdThread models.Thread, err error) {

	query := `INSERT INTO Thread(slug, author, created, forum, message, title) VALUES
						 	($1, $2, $3, $4, $5, $6) 
						 `
	queryAddThreadReturning(&query)
	row := tx.QueryRow(query, thread.Slug, thread.Author, thread.Created,
		thread.Forum, thread.Message, thread.Title)

	createdThread, err = threadScan(row)
	return
}

// updatedThread
func (db *DataBase) threadUpdate(tx *sql.Tx, thread *models.Thread, slug string) (updatedThread models.Thread, err error) {

	query := queryUpdateThread(thread.Message, thread.Title)
	if query == "" {
		updatedThread, err = db.threadFindByIDorSlug(tx, slug)
		return
	}
	queryAddSlug(&query, slug)
	queryAddThreadReturning(&query)
	fmt.Print("threadUpdate query:", query)
	row := tx.QueryRow(query)
	updatedThread, err = threadScan(row)
	return
}

// getThreads get threads
func (db *DataBase) threadsGetWithLimit(tx *sql.Tx, slug string, limit int) (foundThreads []models.Thread, err error) {

	query := querySelectThread() + ` where forum like $1 Limit $2 `

	var rows *sql.Rows

	if rows, err = tx.Query(query, slug, limit); err != nil {
		return
	}
	defer rows.Close()

	foundThreads = []models.Thread{}
	for rows.Next() {
		if err = threadsScan(rows, &foundThreads); err != nil {
			break
		}
	}
	return
}

func (db *DataBase) threadsGet(tx *sql.Tx, slug string, limit int, lb bool, t time.Time, tb bool, desc bool) (foundThreads []models.Thread, err error) {

	query := querySelectThread() + ` where lower(forum) like lower($1)`

	if tb {
		if desc {
			query += ` and created <= $2`
			query += ` order by created desc`
		} else {
			query += ` and created >= $2`
			query += ` order by created`
		}
		if lb {
			query += ` Limit $3`
		}
	} else if lb {
		if desc {
			query += ` order by created desc`
		} else {
			query += ` order by created`
		}
		query += ` Limit $2`
	}

	var rows *sql.Rows

	if tb {
		if lb {
			rows, err = tx.Query(query, slug, t, limit)
		} else {
			rows, err = tx.Query(query, slug, t)
		}
	} else if lb {
		rows, err = tx.Query(query, slug, limit)
	} else {
		rows, err = tx.Query(query, slug)
	}

	if err != nil {
		return
	}
	defer rows.Close()

	foundThreads = []models.Thread{}
	for rows.Next() {
		if err = threadsScan(rows, &foundThreads); err != nil {
			break
		}
	}
	return
}

func (db DataBase) threadConfirmUnique(tx *sql.Tx, thread *models.Thread) (foundThread models.Thread, err error) {
	// if foundThread, err = db.threadFindByTitle(tx, thread.Title); err != sql.ErrNoRows {
	// 	err = re.ErrorThreadConflict()
	// 	return
	// }
	if thread.Slug != "" {
		if foundThread, err = db.threadFindBySlug(tx, thread.Slug); err != sql.ErrNoRows {
			err = re.ErrorThreadConflict()
			return
		}
	}
	err = nil
	return
}

func (db DataBase) threadFindByTitle(tx *sql.Tx, title string) (foundThread models.Thread, err error) {

	query := querySelectThread() + ` where title like $1`

	row := tx.QueryRow(query, title)
	foundThread, err = threadScan(row)
	return
}

func (db DataBase) threadFindBySlug(tx *sql.Tx, slug string) (foundThread models.Thread, err error) {

	query := querySelectThread() + `where lower(slug) like lower($1)`

	row := tx.QueryRow(query, slug)
	foundThread, err = threadScan(row)
	return
}

func (db DataBase) threadFindByID(tx *sql.Tx, arg int) (foundThread models.Thread, err error) {

	query := querySelectThread() + `  where id like $1`

	row := tx.QueryRow(query, arg)
	foundThread, err = threadScan(row)
	return
}

func (db *DataBase) threadCheckID(tx *sql.Tx, oldID int) (newID int, err error) {
	var thatThread models.Thread
	if thatThread, err = db.threadFindByID(tx, oldID); err != nil {
		err = re.ErrorThreadNotExist()
		return
	}
	newID = thatThread.ID
	return
}

func (db DataBase) threadFindByIDorSlug(tx *sql.Tx, arg string) (foundThread models.Thread, err error) {

	query := querySelectThread()
	queryAddSlug(&query, arg)
	row := tx.QueryRow(query)
	foundThread, err = threadScan(row)
	return
}

// addings to query

// queryAddSlug identifier thread by slug_or_id
func queryAddSlug(query *string, arg string) {

	if _, err := strconv.Atoi(arg); err != nil {
		*query += ` where lower(slug) like lower('` + arg + `')`
	} else {
		*query += ` where id = ` + arg
	}
}

// queryAddThreadReturning add returning for insert,update etc
func queryAddThreadReturning(query *string) {
	*query += ` RETURNING id, slug, author, created, forum, message, title, votes `
}

// queryAddThreadReturning add returning for insert,update etc
func querySelectThread() string {
	return ` SELECT T.id, T.slug, T.author, T.created, T.forum, T.message, T.title, T.votes from Thread as T `
}

// scan row to model Vote
func threadScan(row *sql.Row) (foundThread models.Thread, err error) {
	foundThread = models.Thread{}
	err = row.Scan(&foundThread.ID, &foundThread.Slug,
		&foundThread.Author, &foundThread.Created, &foundThread.Forum,
		&foundThread.Message, &foundThread.Title, &foundThread.Votes)
	return
}

// scan rows to model Vote
func threadsScan(rows *sql.Rows, foundThreads *[]models.Thread) (err error) {
	foundThread := models.Thread{}
	err = rows.Scan(&foundThread.ID, &foundThread.Slug,
		&foundThread.Author, &foundThread.Created, &foundThread.Forum,
		&foundThread.Message, &foundThread.Title, &foundThread.Votes)
	if err == nil {
		*foundThreads = append(*foundThreads, foundThread)
	}
	return
}
