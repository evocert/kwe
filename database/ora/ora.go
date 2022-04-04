package ora

import (
	"database/sql"
	"net/url"

	"github.com/evocert/kwe/database"
	//helper registration oracle server driver
	_ "github.com/sijms/go-ora"
	v2 "github.com/sijms/go-ora/v2"
)

//Open -wrap sql.Open("oracle", datasource)
//or
//Open -wrap sql.Open("oracle-ext", datasource)
func Open(oraname, datasource string) (*sql.DB, error) {
	if url, _ := url.ParseRequestURI(datasource); url != nil {
		return sql.Open(oraname, datasource)
	}
	return nil, nil
}

func init() {
	sql.Register("oracle-xe", &v2.OracleDriver{})
	database.GLOBALDBMS().RegisterDriver("oracle", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		db, err = Open("oracle", datasource)
		return
	})
	database.GLOBALDBMS().RegisterDriver("oracle-ext", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		db, err = Open("oracle-ext", datasource)
		return
	})
}
