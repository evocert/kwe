package db2

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	//helper registration db2 server driver
	_ "github.com/ibmdb/go_ibm_db"
)

//Open -wrap sql.Open("mysql", datasource)
func Open(datasource string) (*sql.DB, error) {
	return sql.Open("go_ibm_db", datasource)
}

func init() {
	database.GLOBALDBMS().RegisterDriver("db2", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		db, err = Open(datasource)
		return
	})
}
