package jsext

import (
	"os"
	"path/filepath"
	"runtime"
)

func Register_jsext_osutils(lclobjmp map[string]interface{}) {
	if lclobjmp == nil {
		return
	}
	//vm.SetFieldNameMapper(goja.TagFieldNameMapper("json",true))
	type Version struct {
		Major int `json:"major"`
		Minor int `json:"minor"`
		Bump  int `json:"bump"`
	}
	//todo: namespace everything kwe.fsutils.etcetcetc
	//first test for kwe then do set fsutils on kwe
	lclobjmp["osutils"] = map[string]interface{}{
		"version": Version{
			Major: 0,
			Minor: 0,
			Bump:  2,
		},
		"goos": func() string {
			return runtime.GOOS
		},
		"setEnv": func(name string, val string) {
			os.Setenv(name, val)
		},
		"getEnv": func(name string) string {
			return os.Getenv(name)
		},
		"clearEnv": func() {
			os.Clearenv()
		},
		"env": func() []string {
			return os.Environ()
		},
		"goarch": func() string {
			return runtime.GOARCH
		},
		"pwd": func() string {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				panic("Failed to obtain root dir")
			}
			return dir
		},
	}
}
