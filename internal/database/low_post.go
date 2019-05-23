package database

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/SmartPhoneJava/forum_bd/internal/models"
	//re "github.com/SmartPhoneJava/forum_bd/internal/return_errors"

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

	// var (
	// 	path string
	// )
	// if post.Parent == 0 {
	// 	path = "0"
	// } else {
	// 	if path, err = getPath(tx, post.Parent); err != nil {
	// 		return
	// 	}
	// 	if path == "" {
	// 		err = re.ErrorInvalidPath()
	// 	}
	// }

	query := `INSERT INTO Post(author, created, forum, message, thread, parent) VALUES
						 	($1, $2, $3, $4, $5, $6) 
						 `
	queryAddPostReturning(&query)
	row := tx.QueryRow(query, post.Author, t,
		thread.Forum, post.Message, thread.ID, post.Parent)

	if createdPost, err = postScan(row); err != nil {
		return
	}

	// query = `UPDATE Post set path=$1 where id=$2 `
	// updatePath(&path, createdPost.ID)
	// _, err = tx.Exec(query, path, createdPost.ID)

	return
}

func intToString(n int) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func addPostToQuery(post models.Post) string {
	return "('" +
		post.Author + "',$1,$2,'" +
		post.Message + "', $3,'" +
		intToString(post.Parent) + "')"
}

func addPostAuthor(post models.Post) string {
	return "(lower('" + post.Author + "'),lower($1)"
}

func (db *DataBase) userInForumCreatePosts(posts []models.Post, thread models.Thread) {
	for _, post := range posts {
		db.userInForumCreate(post.Author, thread.Forum)
	}
	return
}

// postCreate create post
func (db *DataBase) postsCreate(tx *sql.Tx, posts []models.Post, thread models.Thread,
	t time.Time) (err error) {
	if len(posts) == 0 {
		return nil
	}

	query := `
		INSERT INTO Post(author, created, forum, message, thread, parent) VALUES
						
						 `

	for i, post := range posts {
		if i == 0 {
			query += addPostToQuery(post)
		} else {
			query += "," + addPostToQuery(post)
		}
		posts[i].Thread = thread.ID
		posts[i].Forum = thread.Forum
	}
	query += " returning id, created"
	//queryAddPostReturning(&query)
	debug("query createPosts:", query)
	debug("query pars:", thread.Forum, thread.ID)
	var rows *sql.Rows

	if rows, err = tx.Query(query, t, thread.Forum, thread.ID); err != nil {
		debug("error is here", err.Error())
		return
	}
	defer rows.Close()
	debug("try")

	//createdPosts = []models.Post{}
	i := 0
	for rows.Next() {
		if err = rows.Scan(&posts[i].ID, &posts[i].Created); err != nil {
			break
		}
		i++
	}
	debug("size:", i)

	if err != nil {
		debug(err.Error())
	}
	debug("done")
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
	foundPosts, err = db.postsGet(tx, thread, qc, pq, 0)
	return
}

// postsGetTree
func (db *DataBase) postsGetTree(tx *sql.Tx,
	thread models.Thread, qc QueryGetConditions) (foundPosts []models.Post, err error) {

	var path string

	if qc.mn {
		if path, err = getPath(tx, qc.mv); err != nil {
			debug("err!!", err.Error())
			return
		}
	}
	debug("ready!!")
	pq := &postQuery{
		sortASC:     ` order by string_to_array(path, '.')::int[], created `,
		sortDESC:    ` order by string_to_array(path, '.')::int[] desc, created desc `,
		compareASC:  ` and string_to_array(path, '.')::int[] > string_to_array('` + path + `', '.')::int[] `,
		compareDESC: ` and string_to_array(path, '.')::int[] < string_to_array('` + path + `', '.')::int[] `,
	}
	foundPosts, err = db.postsGet(tx, thread, qc, pq, 1)
	return
}

// postsGetParentTree
func (db *DataBase) postsGetParentTree(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions) (foundPosts []models.Post, err error) {

	var path string

	if qc.mn {
		if path, err = getPath(tx, qc.mv); err != nil {
			fmt.Print("err!!", err.Error())
			return
		}
	}
	debug("readywww!!")
	pq := &postQuery{
		sortASC:  ` order by string_to_array(path, '.')::int[], created `,
		sortDESC: ` order by split_part(path, '.', 2)::int desc, string_to_array(path, '.')::int[], created `,
		compareASC: ` and string_to_array(path, '.')::int[] > string_to_array('` + path + `', '.')::int[]
			 and split_part(path, '.', 2)::int > split_part('` + path + `', '.', 2)::int `,
		compareDESC: ` and string_to_array(path, '.')::int[] < string_to_array('` + path + `', '.')::int[]
			 and split_part(path, '.', 2)::int < split_part('` + path + `', '.', 2)::int `,
	}

	foundPosts, err = parentTreeGet(tx, thread, qc, pq)
	return
}

// postsGet
func (db *DataBase) postsGet(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions, pq *postQuery, vvv int) (foundPosts []models.Post, err error) {

	var query string
	/*
		if vvv == 1 {
			query =
				"SELECT t.id, t.author, t.created, t.forum, t.message, t.thread, t.parent, t.path, t.isEdited  " +
					"FROM Post as t where t.thread = $1 "
			if qc.mn {
				if qc.desc {
					query += " and split_part(t.path, '.', 2) < split_part((SELECT path FROM Post WHERE id = " + strconv.Itoa(qc.mv) + "), '.', 2) "
				} else {
					query += " and split_part(t.path, '.', 2) > split_part((SELECT path FROM Post WHERE id = " + strconv.Itoa(qc.mv) + "), '.', 2) "
				}
			}

			query += " order by t.path "
			if qc.desc {
				query += "desc "
			}

			if qc.ln {
				query += "limit " + strconv.Itoa(qc.lv) + ";"
			}
		} else {
	*/
	query = querySelectPost() + `  
		 where thread = $1
	`
	queryAddConditions(&query, qc, pq)
	//}

	debug("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", query)
	if vvv == 1 {
		debug("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", query)
	}

	var rows *sql.Rows

	if rows, err = tx.Query(query, thread.ID); err != nil {
		return
	}
	defer rows.Close()

	foundPosts = []models.Post{}
	for rows.Next() {
		if err = postsScan(rows, &foundPosts); err != nil {
			break
		}
	}
	debug("query:"+query, ":::", len(foundPosts))
	//fmt.Println("size:", len(foundPosts))

	return
}

// parentTreeGet
func parentTreeGet(tx *sql.Tx, thread models.Thread,
	qc QueryGetConditions, pq *postQuery) (foundPosts []models.Post, err error) {

	queryInside := `
		select split_part(path, '.', 2)::int as p
			from Post 
				where thread = $1
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
		where thread = $1 
				 	and split_part(path, '.', 2)::int = ANY 
				 	(` + groupBy + `)
	`
	queryAddSort(&query, qc, pq.sortASC, pq.sortDESC)

	var rows *sql.Rows

	if rows, err = tx.Query(query, thread.ID); err != nil {
		return
	}
	defer rows.Close()

	foundPosts = []models.Post{}
	for rows.Next() {
		if err = postsScan(rows, &foundPosts); err != nil {
			break
		}
	}

	debug("debug me qu:", query, "::::::::", len(foundPosts), thread.ID)

	return
}

// postCheckParent
func (db *DataBase) postCheckParent(tx *sql.Tx, post models.Post, thread models.Thread) (err error) {

	query := `
	select 1
		from Post as P
		where id = $1 and thread = $2
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
