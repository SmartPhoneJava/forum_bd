package database

import (
	"database/sql"

	"github.com/SmartPhoneJava/forum_bd/internal/models"

	//
	_ "github.com/lib/pq"
)

// CreateVote handle vote creation
func (db *DataBase) CreateVote(vote models.Vote, slug string) (thread models.Thread, err error) {

	var (
		tx        *sql.Tx
		threadID  int
		prevVoice int
	)
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if thread, err = db.threadFindByIDorSlug(tx, slug); err != nil {
		return
	}
	threadID = thread.ID

	//if threadID, err = db.threadIDBySlug(tx, slug); err != nil {
	//if thread, err = db.threadFindByIDorSlug(tx, slug); err != nil {
	//return
	//}
	//threadID = thread.ID
	//}
	var prevVote models.Vote

	/*
		if prevVote, err = db.voteFindByThreadAndAuthor(tx, threadID, vote.Author); err != nil && err != sql.ErrNoRows {
			return
		}
		if err != nil {
			if vote, err = db.voteCreate(tx, vote, threadID); err != nil {
				return
			}
		} else {

			if vote, prevVoice, err = db.voteUpdate(tx, vote, threadID); err != nil {
				return
			}
		}
	*/
	vote.Thread = threadID
	if err = db.voteCreate(tx, vote); err == nil {
		prevVoice = 0
		err = tx.Commit()
	} else {
		//err = tx.Commit()
		debug("err_create:", err.Error())
		var (
			txx *sql.Tx
		)
		if txx, err = db.Db.Begin(); err != nil {
			return
		}
		defer txx.Rollback()

		if vote, prevVoice, err = db.voteUpdate(txx, vote, threadID); err != nil {
			debug("err_update:", err.Error())
			return
		}
		if err = txx.Commit(); err != nil {
			return
		}

		//vote.Voice = prevVoice
	}
	/*
		if prevVote, err = db.voteFindByThreadAndAuthor(tx, threadID, vote.Author); err != nil && err != sql.ErrNoRows {
			return
		}

		//vote.Print()

		if err != nil {
			prevVoice = 0
			fmt.Println("create")
			if vote, err = db.voteCreate(tx, vote, threadID); err != nil {
				return
			}
		} else {
			prevVoice = prevVote.Voice
			fmt.Println("update")
			if vote, err = db.voteUpdate(tx, vote, threadID); err != nil {
				err = nil
				return
			}
		}
	*/

	newVoice := vote.Voice - prevVoice
	thread.Votes += newVoice
	debug("#"+vote.Author, "newVoice!", thread.Votes, newVoice, vote.Voice, prevVoice, prevVote.Voice)
	// if thread, err = db.voteThread(tx, newVoice, threadID); err != nil {
	// 	return
	// }

	return
}
