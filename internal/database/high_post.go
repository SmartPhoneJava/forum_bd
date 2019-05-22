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

	// size := len(posts)
	// users := make([]string, size)

	/*
		for _, post := range posts {

			// if _, err = db.userCheckID(tx, post.Author); err != nil {
			// 	return
			// }
			fmt.Println("post.Author:", post.Author)

			if post, err = db.postCreate(tx, post, thatThread, t); err != nil {
				return
			}

			// if post.Parent != 0 {
			// 	if err = db.postCheckParent(tx, post, thatThread); err != nil {
			// 		fmt.Println("re.ErrorPostConflict()", thatThread.ID, post.Thread, post.Thread, post.ID)
			// 		//err = re.ErrorPostConflict()
			// 		return
			// 	}
			// }

			createdPosts = append(createdPosts, post)
			count++
		}
	*/
	//errchan := make(chan error)
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
	//done <- err // it is stop for outter functions

	db.userInForumCreatePosts(posts, thatThread)
	done <- err
	//done <- nil

	/*
			errchan := make(chan error)
		//fmt.Println("ready to put")
		errchan <- err
		//fmt.Println("we put")
		defer close(errchan)

		count := len(posts)
		//fmt.Println("init")
		all := &sync.WaitGroup{}
		all.Add(3)
		go db.postsCreate(tx1, posts, thatThread, t, all, errchan)
		go db.forumUpdatePosts(tx2, thatThread.Forum, count, all, errchan)
		go db.statusAddPost(tx3, count, all, errchan)

		//fmt.Println("wait start:")
		all.Wait()
		//fmt.Println("wait over:")
		var ok bool
		if err, ok = <-errchan; ok && (err != nil) {
			//fmt.Println("err:", err.Error())
			return
		}

		err = tx1.Commit()
		err = tx2.Commit()
		err = tx3.Commit()
	*/
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
