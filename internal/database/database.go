package database

import (
	"database/sql"
	"time"

	//
	_ "github.com/lib/pq"
)

// DataBase consists of *sql.DB
// Support methods Login, Register
type DataBase struct {
	Db *sql.DB
}

// QueryParameters qiery
type QueryParameters struct {
	query  string
	thread int
	forum  int
}

// QueryGetConditions query
type QueryGetConditions struct {
	tv   time.Time // time value
	tn   bool      // time need
	mv   int       // min id value
	mn   bool      // min id need
	nv   string    // nickname value
	nn   bool      // nickname need
	lv   int       // limit value
	ln   bool      // limit need
	desc bool      // desc need
}

// InitUser init user
func (qgc *QueryGetConditions) InitUser(
	nn bool, nv string, ln bool, lv int, desc bool) {
	qgc.nv = nv
	qgc.nn = nn
	qgc.lv = lv
	qgc.ln = ln
	qgc.desc = desc
}

// InitPost init post
func (qgc *QueryGetConditions) InitPost(
	mn bool, mv int, ln bool, lv int, desc bool) {
	qgc.mv = mv
	qgc.mn = mn
	qgc.lv = lv
	qgc.ln = ln
	qgc.desc = desc
}

// InitThread init thread
func (qgc *QueryGetConditions) InitThread(tv time.Time,
	tn bool, lv int, ln bool, desc bool) {
	qgc.tv = tv
	qgc.tn = tn
	qgc.lv = lv
	qgc.ln = ln
	qgc.desc = desc
}
