package jsext

import (
	"log"

	"github.com/dop251/goja"
	"github.com/evocert/kwe/iorw"
)

func Register_jsext_consoleutils(vm *goja.Runtime) {
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
		Version Version              `json:"version"`
		Log     func(...interface{}) `json:"log"`
		Warn    func(...interface{}) `json:"warn"`
		Error   func(...interface{}) `json:"error"`
		Debug   func(...interface{}) `json:"debug"`
		Trace   func(...interface{}) `json:"trace"`
	}{
		Version: Version{
			Major: 0,
			Minor: 0,
			Bump:  1,
		},
		//todo: colors
		Log: func(msg ...interface{}) {
			buf := iorw.NewBuffer()
			buf.Print(msg...)
			lgmsg := buf.String()
			buf.Close()
			buf = nil
			log.Println("LOG:   ", lgmsg)
		},
		Warn: func(msg ...interface{}) {
			buf := iorw.NewBuffer()
			buf.Print(msg...)
			lgmsg := buf.String()
			buf = nil
			log.Println("WARN:  ", lgmsg)
		},
		Error: func(msg ...interface{}) {
			buf := iorw.NewBuffer()
			buf.Print(msg...)
			lgmsg := buf.String()
			buf = nil
			log.Println("ERROR: ", lgmsg)
		},
		Debug: func(msg ...interface{}) {
			buf := iorw.NewBuffer()
			buf.Print(msg...)
			lgmsg := buf.String()
			buf = nil
			log.Println("DEBUG: ", lgmsg)
		},
		Trace: func(msg ...interface{}) {
			buf := iorw.NewBuffer()
			buf.Print(msg...)
			lgmsg := buf.String()
			buf = nil
			log.Println("TRACE: ", lgmsg)
		},
	})
}
