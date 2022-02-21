package ora

import (
	"database/sql"
	"net/url"

	"github.com/evocert/kwe/database"
	//helper registration oracle server driver
	_ "github.com/sijms/go-ora"
)

//Open -wrap sql.Open("oracle", datasource)
func Open(datasource string) (*sql.DB, error) {
	//var altdatasource = ""
	/*var tmpdatasource = ""
	if strings.HasPrefix(datasource, "oracle://") {
		tmpdatasource = datasource[len("oracle://"):]

		if strings.Index(tmpdatasource, "@") > -1 {

			fmt.Println(strings.Replace(url.QueryEscape(tmpdatasource[:strings.Index(tmpdatasource, "@")]), "%3A", ":", -1))
		}
	} else {

	}*/
	if url, _ := url.ParseRequestURI(datasource); url != nil {
		return sql.Open("oracle", datasource)
	}
	return nil, nil
}

func init() {
	database.GLOBALDBMS().RegisterDriver("oracle", func(datasource string, a ...interface{}) (db *sql.DB, err error) {
		db, err = Open(datasource)
		return
	})
}
