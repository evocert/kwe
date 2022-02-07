package jsext

import "github.com/dop251/goja"

var localobjmap map[string]interface{} = make(map[string]interface{})

func Register(vm *goja.Runtime) {
	if len(localobjmap) > 0 {
		for lclk, lclv := range localobjmap {
			vm.Set(lclk, lclv)
		}
	}
}

func init() {
	Register_jsext_osutils(localobjmap)
	Register_jsext_time(localobjmap)
	Register_jsext_executils(localobjmap)
	Register_jsext_consoleutils(localobjmap)
	Register_jsext_base62(localobjmap)
}
