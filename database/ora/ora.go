package ora

import (
	"database/sql"
	"net/url"

	"github.com/evocert/kwe/database"
	//helper registration oracle server driver
	_ "github.com/evocert/kwe/go_ora"
	_ "github.com/evocert/kwe/go_ora/go_oraext"
)

//Open -wrap sql.Open("oracle", datasource)
func Open(oraname, datasource string) (*sql.DB, error) {
	if url, _ := url.ParseRequestURI(datasource); url != nil {
		return sql.Open(oraname, datasource)
	}
	return nil, nil
}

func init() {
	database.GLOBALDBMS().RegisterDriver("oracle", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		db, err = Open("kwe:oracle", datasource)
		return
	})
	database.GLOBALDBMS().RegisterDriver("oracle:ext", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		db, err = Open("kwe:oracle:ext", datasource)
		return
	})
}
