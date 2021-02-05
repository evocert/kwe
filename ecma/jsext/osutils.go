package jsext
import (
	"runtime"
	"github.com/dop251/goja"
)
func Register_jsext_osutils(vm *goja.Runtime) {
	//vm.SetFieldNameMapper(goja.TagFieldNameMapper("json",true))
	type Version struct {
		Major int `json:"major"`
		Minor int `json:"minor"`
		Bump  int `json:"bump"`
	}
	//todo: namespace everything kwe.fsutils.etcetcetc
	//first test for kwe then do set fsutils on kwe
	vm.Set("osutils", struct {
		Version         Version                  `json:"version"`
		GOOS            func()string             `json:"GOOS"`
	}{
		Version: Version{
			Major: 0,
			Minor: 0,
			Bump:  2,
		},
		GOOS: func() string {
			return runtime.GOOS
		},
	})
}
