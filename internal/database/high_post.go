package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"time"

	//
	_ "github.com/lib/pq"
)

// UpdatePost handle post creation
func (db *DataBase) UpdatePost(post models.Post, id int) (updatedPost models.Post, err error) {

	var (
		tx *sql.Tx
	)
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if updatedPost, err = db.postUpdate(tx, post, id); err != nil {
		return
	}
	updatedPost.Print()
	err = tx.Commit()
	return
}

// CreatePost handle post creation
func (db *DataBase) CreatePost(posts []models.Post, slug string) (createdPosts []models.Post, err error) {

	var (
		tx         *sql.Tx
		thatThread models.Thread
		count      int
		t          time.Time
	)
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	createdPosts = []models.Post{}

	if thatThread, err = db.threadFindByIDorSlug(tx, slug); err != nil {

		return
	}

	t = time.Now()
	count = 0
	for _, post := range posts {

		if post.Author, err = db.userCheckID(tx, post.Author); err != nil {
			return
		}

		if post.Parent != 0 {
			if err = db.postCheckParent(tx, post, thatThread); err != nil {
				err = re.ErrorPostConflict()
				return
			}
		}

		if post, err = db.postCreate(tx, post, thatThread, t); err != nil {
			return
		}

		createdPosts = append(createdPosts, post)
		count++
	}
	if err = db.forumUpdatePosts(tx, thatThread.Forum, count); err != nil {
		return
	}

	if err = db.statusAddPost(tx, count); err != nil {
		return
	}

	err = tx.Commit()
	return
}

// GetPosts get posts
func (db *DataBase) GetPosts(slug string, qgc QueryGetConditions, sort string) (returnPosts []models.Post, err error) {

	var tx *sql.Tx
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	var thatThread models.Thread

	if thatThread, err = db.threadFindByIDorSlug(tx, slug); err != nil {
		return
	}

	switch sort {
	case "tree":
		returnPosts, err = db.postsGetTree(tx, thatThread, qgc)
	case "parent_tree":
		returnPosts, err = db.postsGetParentTree(tx, thatThread, qgc)
	default:
		returnPosts, err = db.postsGetFlat(tx, thatThread, qgc)
	}

	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
