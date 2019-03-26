package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"fmt"

	//
	_ "github.com/lib/pq"
)

// createUser
func (db *DataBase) createUser(tx *sql.Tx, user *models.User) (createdUser models.User, err error) {

	query := `INSERT INTO UserForum(fullname, nickname, email, about) VALUES
						 	($1, '` + user.Nickname + `', $2, $3) 
						 `
	queryAddUserReturning(&query)
	row := tx.QueryRow(query, user.Fullname, user.Email, user.About)
	createdUser, err = userScan(row)
	return
}

// updateUser
func (db *DataBase) updateUser(tx *sql.Tx, user models.User) (updated models.User, err error) {

	query := queryUpdateUser(user.Fullname, user.Email, user.About)
	if query == "" {
		updated, err = db.findUserByName(tx, user.Nickname)
		return
	}
	query += `	where nickname = $1 	`
	queryAddUserReturning(&query)
	row := tx.QueryRow(query, user.Nickname)

	updated, err = userScan(row)
	return
}

// userCheckID
func (db *DataBase) userCheckID(tx *sql.Tx, oldNickname string) (newNickname string, err error) {
	var thatUser models.User
	if thatUser, err = db.findUserByName(tx, oldNickname); err != nil {
		err = re.ErrorUserNotExist()
		return
	}
	newNickname = thatUser.Nickname
	return
}

// confirmUnique confirm that user.Email and user.Name
// dont use by another Player
func (db DataBase) userConfirmUnique(tx *sql.Tx, user *models.User) (users *[]models.User, err error) {

	var foundUsers *[]models.User
	users = &[]models.User{}

	foundUsers, err = db.isOnlyEmailUnique(tx, user.Email, user.Nickname)

	if err != nil {
		return
	}

	if foundUsers != nil && len(*foundUsers) > 0 {
		*users = append(*users, *foundUsers...)
		foundUsers = nil
	}

	foundUsers, err = db.isNicknameUnique(tx, user.Nickname)

	if err != nil {
		return
	}

	if foundUsers != nil && len(*foundUsers) > 0 {
		*users = append(*users, *foundUsers...)
		foundUsers = nil
	}

	return
}

// findUser
func (db DataBase) findUser(tx *sql.Tx, queryAdd string, arg string) (foundUser models.User, err error) {

	query := querySelectUser() + queryAdd
	row := tx.QueryRow(query, arg)
	foundUser, err = userScan(row)
	return
}

// findUserByName
func (db DataBase) findUserByName(tx *sql.Tx, taken string) (foundUser models.User, err error) {

	query := `where lower(nickname) like lower($1)`
	foundUser, err = db.findUser(tx, query, taken)
	return
}

func (db DataBase) findUserByEmail(tx *sql.Tx, taken string) (foundUser models.User, err error) {

	query := `where lower(email) like lower($1)`
	foundUser, err = db.findUser(tx, query, taken)
	return
}

// usersGet
func (db *DataBase) usersGet(tx *sql.Tx, slug string,
	qc QueryGetConditions) (foundUsers []models.User, err error) {

	pq := &postQuery{
		sortASC:     ` order by lower(nickname) `,
		sortDESC:    ` order by lower(nickname) desc `,
		compareASC:  ` and lower(nickname) > lower('` + qc.nv + `')`,
		compareDESC: ` and lower(nickname) < lower('` + qc.nv + `')`,
	}

	query := querySelectUser() + ` as uf 
		where (
			nickname in 
		(
			select author
				from Post
				where 
				lower(uf.nickname) like lower(author) and
				lower(forum) like lower($1)
		) or
		nickname in 
		(
			select author
				from Thread
				where 
				lower(uf.nickname) like lower(author) and
				lower(forum) like lower($1)
		)
		)
	`
	queryAddConditions(&query, qc, pq)

	fmt.Println("query:" + query)
	var rows *sql.Rows

	if rows, err = tx.Query(query, slug); err != nil {
		return
	}
	defer rows.Close()

	foundUsers = []models.User{}
	for rows.Next() {
		if err = usersScan(rows, &foundUsers); err != nil {
			break
		}
	}
	return
}

// findUsers
func (db DataBase) findUsers(tx *sql.Tx, queryAdd string, taken ...string) (foundUsers *[]models.User, err error) {

	query := querySelectUser() + queryAdd

	var rows *sql.Rows

	if len(taken) == 1 {
		rows, err = tx.Query(query, taken[0])
	} else {
		rows, err = tx.Query(query, taken[0], taken[1])
	}

	if err != nil {
		return
	}
	defer rows.Close()

	foundUsers = &[]models.User{}
	for rows.Next() {
		if err = usersScan(rows, foundUsers); err != nil {
			break
		}
	}
	return
}

// isNicknameUnique checks if there are Players with
// this('taken') nickname and returns corresponding error if yes. A
func (db DataBase) isNicknameUnique(tx *sql.Tx, taken string) (foundUsers *[]models.User, err error) {

	query := `where lower(nickname) like lower($1)`
	foundUsers, err = db.findUsers(tx, query, taken)
	return
}

// isEmailUnique checks if there are Players with
// this('taken') email and returns corresponding error if yes. B\A
func (db DataBase) isOnlyEmailUnique(tx *sql.Tx, email string, nickname string) (foundUsers *[]models.User, err error) {

	query := `where lower(email) like lower($1) and lower(nickname) not like lower($2)`
	foundUsers, err = db.findUsers(tx, query, email, nickname)
	return
}

// query add returning
func queryAddUserReturning(query *string) {
	*query += ` RETURNING fullname, nickname, email, about `
}

// scan row to model User
func userScan(row *sql.Row) (foundUser models.User, err error) {
	foundUser = models.User{}
	err = row.Scan(&foundUser.Fullname, &foundUser.Nickname,
		&foundUser.Email, &foundUser.About)
	return
}

// scan rows to model User
func usersScan(rows *sql.Rows, foundUsers *[]models.User) (err error) {
	foundUser := models.User{}
	err = rows.Scan(&foundUser.Fullname, &foundUser.Nickname,
		&foundUser.Email, &foundUser.About)
	if err == nil {
		*foundUsers = append(*foundUsers, foundUser)
	}
	return
}

// querySelectUser select
func querySelectUser() string {
	return ` SELECT fullname, nickname, email, about FROM UserForum `
}
