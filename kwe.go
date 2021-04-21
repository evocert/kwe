package main

import (
	"os"

	//_ "github.com/evocert/kwe/database/db2"

	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"
	"github.com/evocert/kwe/service"
)

func main() {
	//mphndl := caching.NewMapHandler()
	//mphndl.Put("k1", 89890809)
	service.RunService(os.Args...)
}
