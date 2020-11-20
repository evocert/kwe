package jsext

import (
	"log"

	"github.com/evocert/kwe/ecma/es51"
)

func Register_jsext_consoleutils(vm *es51.Runtime) {
	//vm.SetFieldNameMapper(goja.TagFieldNameMapper("json",true))
	type Version struct {
		Major int `json:"major"`
		Minor int `json:"minor"`
		Bump  int `json:"bump"`
	}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds) //global???
	//todo: namespace everything kwe.fsutils.etcetcetc
	//first test for kwe then do set fsutils on kwe
	vm.Set("console", struct {
		Version Version      `json:"version"`
		Log     func(string) `json:"log"`
		Warn    func(string) `json:"warn"`
		Error   func(string) `json:"error"`
		Debug   func(string) `json:"debug"`
		Trace   func(string) `json:"trace"`
	}{
		Version: Version{
			Major: 0,
			Minor: 0,
			Bump:  1,
		},
		//todo: colors
		Log: func(msg string) {
			log.Println("LOG:   ", msg)
		},
		Warn: func(msg string) {
			log.Println("WARN:  ", msg)
		},
		Error: func(msg string) {
			log.Println("ERROR: ", msg)
		},
		Debug: func(msg string) {
			log.Println("DEBUG: ", msg)
		},
		Trace: func(msg string) {
			log.Println("TRACE: ", msg)
		},
	})
}
