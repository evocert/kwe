//go:build database && all
// +build database,all

package imports

import (
	_ "github.com/evocert/kwe/database/kwesqlite"
	_ "github.com/evocert/kwe/database/mysql"
	_ "github.com/evocert/kwe/database/ora"
	_ "github.com/evocert/kwe/database/sqlite"
	_ "github.com/evocert/kwe/database/sqlserver"
)
