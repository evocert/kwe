package jsext
import (
	"runtime"
	"github.com/dop251/goja"
	"os"
	"path/filepath"	
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
		Pwd             func()string             `json:"pwd"`
	}{
		Version: Version{
			Major: 0,
			Minor: 0,
			Bump:  2,
		},
		GOOS: func() string {
			return runtime.GOOS
		},
		Pwd: func() string {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err!=nil{
				panic(vm.ToValue("Failed to obtain root dir"))
			}
			return dir;
		},
	})
}
