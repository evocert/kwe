package postgres

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func init() {
	database.GLOBALDBMS().RegisterDriver("postgres", func(datasource string) (db *sql.DB, err error) {
		db, err = sql.Open("pgx", datasource)
		return
	})
}
