package database

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/SmartPhoneJava/forum_bd/internal/models"
	re "github.com/SmartPhoneJava/forum_bd/internal/return_errors"

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

	debug("query", query)
	debug("pars", thread.Slug, thread.Author, thread.Created,
		thread.Forum, thread.Message, thread.Title)
	createdThread, err = threadScan(row)
	/*
		var id string
		queryID := `select nextval('thread_id_seq');`
		row := tx.QueryRow(queryID)
		if err = row.Scan(&id); err != nil {
			return
		}

		query := `INSERT INTO Thread(id, slug, author, created, forum, message, title) VALUES
							 	($1, $2, $3, $4, $5, $6, $7)
							 `
		//queryAddThreadReturning(&query)
		_, err = tx.Exec(query, id, thread.Slug, thread.Author, thread.Created,
			thread.Forum, thread.Message, thread.Title)

		debug("query", "INSERT INTO Thread(id, slug, author, forum, message, title) VALUES ("+
			id+",'"+thread.Slug+"', '"+thread.Author+"', '", thread.Forum, "','"+
			thread.Message+"', '"+thread.Title+"') ")
		debug("pars", thread.Slug, thread.Author, thread.Created,
			thread.Forum, thread.Message, thread.Title)
		debug("thread id:", id)

		if err != nil {
			debug("cant create cause:", err.Error())
			return
		}

		queryGet := `select id, slug, author, created, forum, message, title, votes from Thread where id = ` + id
		row = tx.QueryRow(queryGet)
		createdThread, err = threadScan(row)
		// if err == sql.ErrNoRows {
		// 	err = nil
		// 	return
		// }
	*/
	if err != nil {
		debug("err create:", err.Error())
	}
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
	debug("threadUpdate query:", query)
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

func (db DataBase) threadIDBySlug(tx *sql.Tx, slug string) (id int, err error) {

	query := ` SELECT id from Thread`
	if id, err = strconv.Atoi(slug); err != nil {
		query += ` where lower(slug) like lower($1)`
		row := tx.QueryRow(query, slug)
		err = row.Scan(&id)
	} else {
		return id, nil
	}
	debug("slug", slug)
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
