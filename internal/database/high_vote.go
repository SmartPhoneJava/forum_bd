package database

import (
	"database/sql"
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

// CreateVote handle vote creation
func (db *DataBase) CreateVote(vote models.Vote, slug string) (thread models.Thread, err error) {

	var (
		tx        *sql.Tx
		prevVote  models.Vote
		prevVoice int
	)
	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if thread, err = db.threadFindByIDorSlug(tx, slug); err != nil {
		return
	}

	if prevVote, err = db.voteFindByThreadAndAuthor(tx, thread.ID, vote.Author); err != nil && err != sql.ErrNoRows {
		return
	}

	//vote.Print()

	if err != nil {
		prevVoice = 0
		if vote, err = db.voteCreate(tx, vote, thread); err != nil {
			return
		}
	} else {
		prevVoice = prevVote.Voice
		if vote, err = db.voteUpdate(tx, vote, thread); err != nil {
			err = nil
			return
		}
	}

	newVoice := vote.Voice - prevVoice
	if thread, err = db.voteThread(tx, newVoice, thread); err != nil {
		return
	}
	err = tx.Commit()
	return
}
