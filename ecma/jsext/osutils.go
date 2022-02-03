package jsext

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	base62 "github.com/evocert/kwe/encoding/base62"
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
		"args": func() []string {
			return os.Args[1:]
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

func Register_jsext_time(lclobjmp map[string]interface{}) {
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
	lclobjmp["timing"] = map[string]interface{}{
		"version": Version{
			Major: 0,
			Minor: 0,
			Bump:  1,
		},
		"nanoToday": func() (nanotoday uint64) {
			t := time.Now()
			nanotoday = uint64(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).UnixNano())
			return
		},
		"nanoMidnight": func() (nanomidnight uint64) {
			t := time.Now()
			nanomidnight = uint64(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1).Add(time.Nanosecond - 1).UnixNano())
			return
		},
		"nanoNow": func() (nanomidnight uint64) {
			time.Sleep(time.Nanosecond * 1)
			nanomidnight = uint64(time.Now().UnixNano())
			return
		},
	}
}

func Register_jsext_base62(lclobjmp map[string]interface{}) {
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
	lclobjmp["base62"] = map[string]interface{}{
		"version": Version{
			Major: 0,
			Minor: 0,
			Bump:  1,
		},
		"encode": func(number uint64) string {
			return base62.Encode(number)
		},
		"decode": func(encoded string) (number uint64) {
			number, _ = base62.Decode(encoded)
			return
		},
	}
}
