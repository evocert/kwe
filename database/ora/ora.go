package ora

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	//helper registration mysql server driver
	_ "github.com/sijms/go-ora"
)

//Open -wrap sql.Open("mysql", datasource)
func Open(datasource string) (*sql.DB, error) {
	return sql.Open("oracle", datasource)
}

func init() {
	database.GLOBALDBMS().RegisterDriver("oracle", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		db, err = Open(datasource)
		return
	})
}
