package jsext

import "github.com/dop251/goja"

func Register(vm *goja.Runtime) {
	Register_jsext_fsutils(vm)
	Register_jsext_osutils(vm)
	Register_jsext_httputils(vm)
	Register_jsext_gfxutils(vm)
	Register_jsext_executils(vm)
	Register_jsext_consoleutils(vm)
}
