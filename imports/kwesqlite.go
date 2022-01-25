//go:build database && kwesqlite
// +build database,kwesqlite

package imports

import (
	_ "github.com/evocert/kwe/database/kwesqlite"
)