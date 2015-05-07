/*
This is a cache that creates, populates and reads entries to an SQLite
database. 

This package is intended to be used as a series of goroutines. They
communicate by channels. The DefinitionWriter will create a new SQLite
database if one does not exist. 

Author: Justin Cook <jhcook@gmail.com>
*/

package cache_sqlite

import (
    "github.com/jhcook/game_engine/util"
	"database/sql"
    "fmt"
    "log"
    "time"
	_ "github.com/mattn/go-sqlite3"
)

type WordDefinition struct {
	uid        sql.NullInt64
	Word       sql.NullString
	Definition sql.NullString
	tcreated   sql.NullInt64
}

var DB_DRIVER string = "sqlite3"

// The following is a side effect of importing go-sqlite3
//sql.Register(DB_DRIVER, &sqlite3.SQLiteDriver{})

func DefinitionWriter(inDefs chan []string, reqs chan string,
                      outDefs chan *WordDefinition) {
	dbw, cerr := sql.Open(DB_DRIVER, "/tmp/sqlite.db")
	if cerr != nil {
		panic(cerr)
	}

	tx, cerr := dbw.Begin()
	if cerr != nil {
		panic(cerr)
	}

	_, cerr = dbw.Exec(`CREATE TABLE IF NOT EXISTS Words (
                        uid INTEGER PRIMARY KEY AUTOINCREMENT,
                        word varchar(32) NOT NULL,
                        definition varchar(512) NOT NULL,
                        tcreated INTEGER(10))`)
	if cerr != nil {
		panic(cerr)
	}
	tx.Commit()

	for {
		iStf := <-inDefs
        // Check to see if already in the db
        reqs <- iStf[0]
        wrd := <- outDefs
        if wrd != nil {
            log.Println(fmt.Sprintf("%s: %s exists", util.GetFuncName(), wrd.Word.String))
            continue
        }

        log.Println(fmt.Sprintf("%s: inserting %s", util.GetFuncName(), iStf[0]))
		tx, cerr = dbw.Begin()
		if cerr != nil {
			panic(cerr)
		}
		stmt, cerr := dbw.Prepare(`INSERT INTO Words (word, definition, tcreated) 
                                  VALUES (?, ?, ?)`)
		if cerr != nil {
			panic(cerr)
		}

		_, cerr = stmt.Exec(iStf[0], iStf[1], time.Now().Unix())
		if cerr != nil {
			panic(cerr)
		}
		tx.Commit()
	}
}

func DefinitionReader(reqs chan string, oDefs chan *WordDefinition) {
	dbr, cerr := sql.Open(DB_DRIVER, "/tmp/sqlite.db")
	if cerr != nil {
		panic(cerr)
	}

	for {
		req := <-reqs
		rows, cerr := dbr.Query("SELECT * FROM Words WHERE Word=? LIMIT 1", req)
		if cerr != nil {
			panic(cerr)
		}

        /* Since it is difficult to get the len(rows) to determine if there are
         * results, make a result array to populate so len(result) will perform
         * that function.
         */
        result := make([]WordDefinition, 0, 2)

		for rows.Next() {
			var uid, tcreated sql.NullInt64
			var word, definition sql.NullString

			if err := rows.Scan(&uid, &word, &definition, &tcreated); err != nil {
				log.Println(fmt.Sprintf("%s: %q", util.GetFuncName(), err))
				continue
			}
            result = append(result, WordDefinition{uid, word, definition, tcreated})
		}

        if len(result) > 0 {
		    oDefs <- &result[0]
        } else {
            oDefs <- nil
        }

		rows.Close()
	}
}

