package main

import (
	"os"

	//_ "github.com/evocert/kwe/database/db2"

	_ "github.com/evocert/kwe/alertify"
	_ "github.com/evocert/kwe/fonts/material"

	_ "github.com/evocert/kwe/typescript"

	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/postgres"
	_ "github.com/evocert/kwe/database/sqlserver"
	"github.com/evocert/kwe/service"
)

func main() {
	service.RunService(os.Args...)
}
