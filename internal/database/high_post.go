package database

import (
	"database/sql"
	"time"

	"github.com/SmartPhoneJava/forum_bd/internal/models"

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
func (db *DataBase) CreatePost(posts []models.Post, slug string, t time.Time, done chan error) {

	var (
		tx         *sql.Tx
		thatThread models.Thread
		count      int
		err        error
	)
	if tx, err = db.Db.Begin(); err != nil {
		done <- err
		return
	}
	defer tx.Rollback()

	if thatThread, err = db.threadFindByIDorSlug(tx, slug); err != nil {
		done <- err
		return
	}

	count = len(posts)
	if err = db.postsCreate(tx, posts, thatThread, t); err != nil {
		done <- err
		return
	}
	debug("posts created")

	if err = db.forumUpdatePosts(tx, thatThread.Forum, count); err != nil {
		return
	}

	if err = db.statusAddPost(tx, count); err != nil {
		return
	}

	err = tx.Commit()
	done <- err // it is stop for outter functions

	db.userInForumCreatePosts(posts, thatThread)
	//done <- err
	//done <- nil
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
	debug("GetPosts begin")
	if thatThread, err = db.threadFindByIDorSlug(tx, slug); err != nil {
		return
	}
	debug("GetPosts end", sort)
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
