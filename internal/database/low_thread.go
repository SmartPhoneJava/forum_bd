package database

import (
	"database/sql"
	"strconv"

	"github.com/SmartPhoneJava/forum_bd/internal/models"
	re "github.com/SmartPhoneJava/forum_bd/internal/return_errors"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) userInForumCreate(nickname, forum string) (err error) {

	query := `
		INSERT into UserInForum(nickname, forum) values (lower($1), lower($2))
						 `

	if _, err = db.Db.Exec(query, nickname, forum); err != nil {
		debug("error createUserInForum is here", err.Error())
		return
	}

	debug("done")
	return
}

// createThread create thread
func (db *DataBase) threadCreate(tx *sql.Tx, thread *models.Thread) (createdThread models.Thread, err error) {
	var (
		row   *sql.Row
		query string
	)
	if thread.Created != "" {
		query = `INSERT INTO Thread(slug, author, created, forum, message, title) VALUES
		($1, $2, $3, $4, $5, $6) 
		RETURNING id, slug, author, created, forum, message, title, votes
	`
		row = tx.QueryRow(query, thread.Slug, thread.Author, thread.Created,
			thread.Forum, thread.Message, thread.Title)

		foundThread := models.Thread{}
		err = row.Scan(&foundThread.ID, &foundThread.Slug,
			&foundThread.Author, &foundThread.Created, &foundThread.Forum,
			&foundThread.Message, &foundThread.Title, &foundThread.Votes)
		createdThread = foundThread
		debug("thread.Created:", thread.Created)
	} else {
		query = `INSERT INTO Thread(slug, author, forum, message, title) VALUES
		($1, $2, $3, $4, $5) 
		RETURNING id, slug, author, forum, message, title, votes
	`
		row = tx.QueryRow(query, thread.Slug, thread.Author,
			thread.Forum, thread.Message, thread.Title)

		foundThread := models.Thread{}
		err = row.Scan(&foundThread.ID, &foundThread.Slug,
			&foundThread.Author, &foundThread.Forum,
			&foundThread.Message, &foundThread.Title, &foundThread.Votes)
		createdThread = foundThread
	}

	//	debug("query", query)
	//debug("pars", thread.Slug, thread.Author, thread.Created,
	//		thread.Forum, thread.Message, thread.Title)
	//createdThread, err = threadScan(row)

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

func (db *DataBase) threadsGet(tx *sql.Tx, slug string, limit int, lb bool, t string, tb bool, desc bool) (foundThreads []models.Thread, err error) {

	//t = t.Add(4 * time.Hour)
	query := querySelectThread() + ` where lower(forum) like lower($1)`

	if tb {
		debug("T.created ", t)
		if desc {
			query += ` and T.created <= $2`
			query += ` order by T.created desc`
		} else {
			query += ` and T.created >= $2`
			query += ` order by T.created`
		}
		if lb {
			query += ` Limit $3`
		}
	} else if lb {

		if desc {
			query += ` order by T.created desc`
		} else {
			query += ` order by T.created`
		}

		query += ` Limit $2`
	} else {
		if desc {
			query += ` order by T.created desc`
		} else {
			query += ` order by T.created`
		}
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
		if foundThreads[len(foundThreads)-1].ID == 9898 || foundThreads[len(foundThreads)-1].ID == 5186 {
			debug("query", query)
			debug("here", tb, foundThreads[len(foundThreads)-1].ID, slug, t, limit)
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
