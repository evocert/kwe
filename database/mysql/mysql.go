package mysql

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	//helper registration sql server driver
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	database.GLOBALDBMS().RegisterDriver("mysql", func(datasource string) (db *sql.DB, err error) {
		db, err = sql.Open("mysql", datasource)
		return
	})
}
