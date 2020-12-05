package mysql

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	//helper registration mysql server driver
	_ "github.com/go-sql-driver/mysql"
)

//Open -wrap sql.Open("mysql", datasource)
func Open(datasource string) (*sql.DB, error) {
	return sql.Open("mysql", datasource)
}

func init() {
	database.GLOBALDBMS().RegisterDriver("mysql", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		db, err = Open(datasource)
		return
	})
}
