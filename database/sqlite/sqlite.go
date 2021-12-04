package sqlite

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	//helper registration sqlite driver

	_ "github.com/mattn/go-sqlite3"
)

//Open -wrap sql.Open("sqlite", datasource)
func Open(datasource string) (*sql.DB, error) {
	return sql.Open("sqlite3", datasource)
}

func init() {
	database.GLOBALDBMS().RegisterDriver("sqlite", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		db, err = Open(datasource)
		return
	})
}
