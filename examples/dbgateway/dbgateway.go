package main

import (
	"database/sql"

	"github.com/evocert/kwe/database"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"
)

func main() {
	database.GLOBALDBMS().RegisterDriver("bla", myconnect)
	database.GLOBALDBMS().RegisterConnection("test", "sqlserver", "server=localhost\\SQLXPRESS;user id=bcoring;password=bc@r1ng;")
	//database.GLOBALDBMS().CN("test")
}

func myconnect(datasource string) (db *sql.DB, err error) {
	db, err = sql.Open("pqx", datasource)
	return
}
