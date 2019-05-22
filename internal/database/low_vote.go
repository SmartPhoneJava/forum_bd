package database

import (
	"database/sql"

	"github.com/SmartPhoneJava/forum_bd/internal/models"
	re "github.com/SmartPhoneJava/forum_bd/internal/return_errors"

	//
	_ "github.com/lib/pq"
)

// voteCreate
func (db *DataBase) voteCreate(tx *sql.Tx, vote models.Vote) (updatedVote models.Vote, prevVoice int, err error) {

	query := `INSERT INTO Vote(author, voice, thread) VALUES
							 ($1, $2, $3) on conflict(author, thread)  do
							 update set voice = $2
							 RETURNING id, author, voice, thread, isEdited, old_voice 
						 `
	row := tx.QueryRow(query, vote.Author, vote.Voice, vote.Thread)

	updatedVote = models.Vote{}
	var id int
	err = row.Scan(&id, &updatedVote.Author, &updatedVote.Voice,
		&updatedVote.Thread, &updatedVote.IsEdited, &prevVoice)
	if id == 0 {
		err = re.ErrorVoteInvalidAuthor()
	}
	debug("id is", id)
	if err != nil {
		debug("and error:", err.Error())
	}
	return
}

// voteFindByThreadAndAuthor
func (db DataBase) voteFindByThreadAndAuthor(tx *sql.Tx, thread int, author string) (foundVote models.Vote, err error) {

	query := `SELECT author, voice, thread, isEdited FROM Vote where thread = $1 and author = $2`

	row := tx.QueryRow(query, thread, author)
	foundVote, err = voteScan(row)
	return
}

// voteUpdate
func (db DataBase) voteUpdate(tx *sql.Tx, vote models.Vote, threadID int) (updatedVote models.Vote, prevVoice int, err error) {

	query := `	UPDATE Vote set old_voice = voice, voice = $1    --, isEdited = true
		where author = $2 and thread = $3 --and isEdited = false
		RETURNING author, voice, thread, isEdited, old_voice;
	`

	row := tx.QueryRow(query, vote.Voice, vote.Author, threadID)

	updatedVote = models.Vote{}
	err = row.Scan(&updatedVote.Author, &updatedVote.Voice,
		&updatedVote.Thread, &updatedVote.IsEdited, &prevVoice)
	return
}

// voteThread
func (db *DataBase) voteThread(tx *sql.Tx, voice int, threadID int) (updated models.Thread, err error) {

	query := `	UPDATE Thread set votes = votes + $1
								where id = $2
						 `
	queryAddThreadReturning(&query)

	row := tx.QueryRow(query, voice, threadID)

	updated, err = threadScan(row)
	return
}

// query add returning
func queryAddVoteReturning(query *string) {
	*query += ` RETURNING author, voice, thread, isEdited;`
}

// scan to model Vote
func voteScan(row *sql.Row) (foundVote models.Vote, err error) {
	foundVote = models.Vote{}
	err = row.Scan(&foundVote.Author, &foundVote.Voice,
		&foundVote.Thread, &foundVote.IsEdited)
	return
}
