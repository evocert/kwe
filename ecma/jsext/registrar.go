package jsext

import (
	"github.com/evocert/kwe/ecma/es51"
)

func Register(vm *es51.Runtime) {
	Register_jsext_fsutils(vm)
	Register_jsext_httputils(vm)
	Register_jsext_gfxutils(vm)
	Register_jsext_executils(vm)
	Register_jsext_consoleutils(vm)
}
