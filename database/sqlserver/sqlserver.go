package sqlserver

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	//helper registration sql server driver
	_ "github.com/denisenkom/go-mssqldb"
)

//Open -wrap sql.Open("sqlserver", datasource)
func Open(datasource string) (*sql.DB, error) {
	return sql.Open("sqlserver", datasource)
}

func init() {
	database.GLOBALDBMS().RegisterDriver("sqlserver", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		db, err = Open(datasource)
		return
	})
}
