// +build database,all

package main

import (
	_ "github.com/evocert/kwe/database/sqlite"
	_ "github.com/evocert/kwe/database/kwesqlite"
	_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/sqlserver"
	_ "github.com/evocert/kwe/database/mysql"
)