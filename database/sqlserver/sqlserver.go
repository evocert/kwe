package sqlserver

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	//helper registration sql server driver
	_ "github.com/denisenkom/go-mssqldb"
)

func init() {
	database.GLOBALDBMS().RegisterDriver("sqlserver", func(datasource string) (db *sql.DB, err error) {
		db, err = sql.Open("sqlserver", datasource)
		return
	})
}
