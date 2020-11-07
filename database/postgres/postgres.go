package postgres

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	//helper registration posgres server pgx driver
	_ "github.com/jackc/pgx/v4/stdlib"
)

//Open -wrap sql.Open("pgx", datasource)
func Open(datasource string) (*sql.DB, error) {
	return sql.Open("pgx", datasource)
}

func init() {
	database.GLOBALDBMS().RegisterDriver("postgres", func(datasource string) (db *sql.DB, err error) {
		db, err = Open(datasource)
		return
	})
}
