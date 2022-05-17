package contexting

import (
	"github.com/evocert/kwe/scripting/api"
)

type Context struct {
	binding api.Binding
	engine  api.Engine
}
