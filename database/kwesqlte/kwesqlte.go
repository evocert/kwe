package kwesqlte

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	//helper registration sqlite driver

	_ "modernc.org/sqlite"
)

//Open -wrap sql.Open("sqlite", datasource)
func Open(datasource string) (*sql.DB, error) {
	return sql.Open("sqlite", datasource)
}

func init() {
	database.GLOBALDBMS().RegisterDriver("kwesqlite", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		db, err = Open(datasource)
		return
	})
}
