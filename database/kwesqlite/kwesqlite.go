package kwesqlite

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	//helper registration sqlite driver

	_ "modernc.org/sqlite"
)

//Open -wrap sql.Open("sqlite", datasource)
// when registering driver "kwesqlite"
func Open(datasource string) (*sql.DB, error) {
	return sql.Open("sqlite", datasource)
}

func init() {
	database.GLOBALDBMS().RegisterDriver("kwesqlite", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		if datasource == ":memory:" {
			datasource = "file::memory:?mode=memory"
		}
		db, err = Open(datasource)
		return
	})
}
