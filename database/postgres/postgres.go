package postgres

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	//helper registration posgres server pgx driver
	_ "github.com/jackc/pgx/v4/stdlib"
	//_ "github.com/lib/pq"
)

//Open -wrap sql.Open("pgx", datasource)
func Open(datasource string) (*sql.DB, error) {
	return sql.Open("pgx", datasource)
}

func init() {
	database.GLOBALDBMS().RegisterDriver("postgres", func(datasource string, a ...interface{}) (db interface{}, err error) {
		db, err = Open(datasource)
		return
	})
}
