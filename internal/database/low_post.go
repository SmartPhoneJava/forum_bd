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

type postQuery struct {
	sortASC     string
	sortDESC    string
	compareASC  string
	compareDESC string
}

// getPath
func getPath(tx *sql.Tx, id int) (path string, err error) {
	query := `select path from Post where id = $1
						 `

	if err = tx.QueryRow(query, id).Scan(&path); err != nil {
		return
	}
	return
}

// updatePath
func updatePath(path *string, id int) {
	*path = *path + "." + strconv.Itoa(id)
}

// postCreate create post
func (db *DataBase) postCreate(tx *sql.Tx, post models.Post, thread models.Thread,
	t time.Time) (createdPost models.Post, err error) {

	var (
		path string
	)
	if post.Parent == 0 {
		path = "0"
	} else {
		if path, err = getPath(tx, post.Parent); err != nil {
			return
		}
		if path == "" {
			err = re.ErrorInvalidPath()
		}
	}

	query := `INSERT INTO Post(author, created, forum, message, thread, parent, path) VALUES
						 	($1, $2, $3, $4, $5, $6, $7) 
						 `
	queryAddPostReturning(&query)
	row := tx.QueryRow(query, post.Author, t,
		thread.Forum, post.Message, thread.ID, post.Parent, path)

	if createdPost, err = postScan(row); err != nil {
		return
	}

	query = `UPDATE Post set path=$1 where id=$2 `
	updatePath(&path, createdPost.ID)
	_, err = tx.Exec(query, path, createdPost.ID)

	return
}

// postFind
func (db *DataBase) postFind(tx *sql.Tx, id int) (foundPost models.Post, err error) {
	query := querySelectPost() + ` where id=$1 `
	foundPost, err = postScan(tx.QueryRow(query, id))
	return
}

// postCreate create post
func (db *DataBase) postUpdate(tx *sql.Tx, post models.Post, id int) (updatedPost models.Post, err error) {

	if updatedPost, err = db.postFind(tx, id); err != nil {
		return
	}

	if updatedPost.Message == post.Message {
		return
	}
	query := queryUpdatePost(post.Message)
	if query == "" {
		return
	}
	query += ` where id=$1 `
	queryAddPostReturning(&query)
	updatedPost, err = postScan(tx.QueryRow(query, id))

	return
}

// queryAddConditions
func queryAddConditions(query *string, qc QueryGetConditions, pq *postQuery) {
	queryAddMinID(query, qc, pq.compareASC, pq.compareDESC)
	queryAddNickname(query, qc, pq.compareASC, pq.compareDESC)
	queryAddSort(query, qc, pq.sortASC, pq.sortDESC)
	queryAddLimit(query, qc)
}

// queryAddSort
func queryAddSort(query *string, qc QueryGetConditions, sortASC string, sortDESC string) {
	if qc.desc {
		*query += sortDESC
	} else {
		*query += sortASC
	}
}

// queryAddNickname
func queryAddNickname(query *string, qc QueryGetConditions, compareIDASC string, compareIDDESC string) {
	if qc.nn {
		if qc.desc {
			*query += compareIDDESC
		} else {
			*query += compareIDASC
		}
	}
}

// queryAddMinID
func queryAddMinID(query *string, qc QueryGetConditions, compareIDASC string, compareIDDESC string) {
	if qc.mn {
		if qc.desc {
			*query += compareIDDESC
		} else {
			*query += compareIDASC
		}
	}
}

// queryAddLimit
func queryAddLimit(query *string, qc QueryGetConditions) {
	if qc.ln {
		*query += ` Limit ` + strconv.Itoa(qc.lv)
	}
	return
}

// postsGetFlat
func (db *DataBase) postsGetFlat(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions) (foundPosts []models.Post, err error) {

	pq := &postQuery{
		sortASC:     ` order by created, id `,
		sortDESC:    ` order by created desc, id desc `,
		compareASC:  `and id > ` + strconv.Itoa(qc.mv),
		compareDESC: `and id < ` + strconv.Itoa(qc.mv),
	}
	foundPosts, err = db.postsGet(tx, thread, qc, pq)
	return
}

// postsGetTree
func (db *DataBase) postsGetTree(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions) (foundPosts []models.Post, err error) {

	var path string

	if qc.mn {
		if path, err = getPath(tx, qc.mv); err != nil {
			return
		}
	}
	pq := &postQuery{
		sortASC: ` order by string_to_array(path, '.')::int[], created `,
		sortDESC: ` order by	string_to_array(path, '.')::int[] desc, created desc `,
		compareASC:  ` and path > '` + path + `'`,
		compareDESC: ` and path < '` + path + `'`,
	}
	foundPosts, err = db.postsGet(tx, thread, qc, pq)
	return
}

// postsGetParentTree
func (db *DataBase) postsGetParentTree(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions) (foundPosts []models.Post, err error) {

	var path string

	if qc.mn {
		if path, err = getPath(tx, qc.mv); err != nil {
			return
		}
	}
	pq := &postQuery{
		sortASC:    ` order by string_to_array(path, '.')::int[], created `,
		sortDESC:   ` order by split_part(path, '.', 2) desc, string_to_array(path, '.')::int[], created `,
		compareASC: ` and path > '` + path + `'`,
		compareDESC: ` and path < '` + path + `
			' and split_part(path, '.', 2) < split_part('` + path + `', '.', 2)`,
	}

	foundPosts, err = parentTreeGet(tx, thread, qc, pq)
	return
}

// postsGet
func (db *DataBase) postsGet(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions, pq *postQuery) (foundPosts []models.Post, err error) {

	query := querySelectPost() + `  
		 where thread = $1 and 
			 lower(forum) like lower($2)
	`
	queryAddConditions(&query, qc, pq)

	fmt.Println("query:" + query)
	var rows *sql.Rows

	if rows, err = tx.Query(query, thread.ID, thread.Forum); err != nil {
		return
	}
	defer rows.Close()

	foundPosts = []models.Post{}
	for rows.Next() {
		if err = postsScan(rows, &foundPosts); err != nil {
			break
		}
	}

	return
}

// parentTreeGet
func parentTreeGet(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions, pq *postQuery) (foundPosts []models.Post, err error) {

	queryInside := `
		select split_part(path, '.', 2) as p
			from Post 
				where thread = $1 and lower(forum) like lower($2)
	`
	queryAddMinID(&queryInside, qc, pq.compareASC, pq.compareDESC)

	groupBy := `
		select A.p from 
		( 
			` + queryInside + ` 
		) as A
		GROUP BY A.p
	`
	queryAddSort(&groupBy, qc, "order by A.p", "order by A.p desc")
	queryAddLimit(&groupBy, qc)

	query := querySelectPost() + ` 
		where thread = $1 and lower(forum) like lower($2)
				 	and split_part(path, '.', 2) = ANY 
				 	(` + groupBy + `)
	`
	queryAddSort(&query, qc, pq.sortASC, pq.sortDESC)

	var rows *sql.Rows

	if rows, err = tx.Query(query, thread.ID, thread.Forum); err != nil {
		return
	}
	defer rows.Close()

	foundPosts = []models.Post{}
	for rows.Next() {
		if err = postsScan(rows, &foundPosts); err != nil {
			break
		}
	}

	return
}

// postCheckParent
func (db *DataBase) postCheckParent(tx *sql.Tx, post models.Post, thread models.Thread) (err error) {

	query := `
	select 1
		from Post as P
		where id!=$1 and $2  =
			(
				select thread from Post where id=$1
			)	
	`

	var tmp int
	err = tx.QueryRow(query, post.Parent, thread.ID).Scan(&tmp)
	return
}

// querySelectPost
func querySelectPost() string {
	return ` SELECT id, author, created, forum,
	 message, thread, parent, path, isEdited FROM Post `
}

// queryAddPostReturning
func queryAddPostReturning(query *string) {
	*query += ` RETURNING id, author, created,
	 forum, message, thread, parent, path, isEdited `
}

// postScan scan row to model Vote
func postScan(row *sql.Row) (foundPost models.Post, err error) {
	foundPost = models.Post{}
	err = row.Scan(&foundPost.ID, &foundPost.Author, &foundPost.Created,
		&foundPost.Forum, &foundPost.Message, &foundPost.Thread, &foundPost.Parent,
		&foundPost.Path, &foundPost.IsEdited)
	return
}

// postScan scan row to model Vote
func postsScan(rows *sql.Rows, foundPosts *[]models.Post) (err error) {
	foundPost := models.Post{}
	err = rows.Scan(&foundPost.ID, &foundPost.Author, &foundPost.Created,
		&foundPost.Forum, &foundPost.Message, &foundPost.Thread, &foundPost.Parent,
		&foundPost.Path, &foundPost.IsEdited)
	if err == nil {
		*foundPosts = append(*foundPosts, foundPost)
	}
	return
}

// 280 -> 307 -> 344
