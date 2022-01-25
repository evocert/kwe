// +build database,ora

package imports

import (
	//To use ora import use go 1.6+
	_ "github.com/evocert/kwe/database/ora"
)