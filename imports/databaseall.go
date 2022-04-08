//go:build databaseall
// +build databaseall

package imports

import (
	_ "github.com/evocert/kwe/database/kwesqlite"
	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/sqlite"
	_ "github.com/evocert/kwe/database/sqlserver"
)
