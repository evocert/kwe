package jsext

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

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
			//log.Println("LOG:   ", lgmsg)
			logger.Output(2, fmt.Sprintln("LOG:   "+lgmsg))
		},
		Warn: func(msg ...interface{}) {
			buf := iorw.NewBuffer()
			buf.Print(msg...)
			lgmsg := buf.String()
			buf = nil
			//log.Println("WARN:  ", lgmsg)
			logger.Output(2, fmt.Sprintln("WARN:   "+lgmsg))
		},
		Error: func(msg ...interface{}) {
			buf := iorw.NewBuffer()
			buf.Print(msg...)
			lgmsg := buf.String()
			buf = nil
			//log.Println("ERROR: ", lgmsg)
			logger.Output(2, fmt.Sprintln("ERROR:   "+lgmsg))
		},
		Debug: func(msg ...interface{}) {
			buf := iorw.NewBuffer()
			buf.Print(msg...)
			lgmsg := buf.String()
			buf = nil
			//log.Println("DEBUG: ", lgmsg)
			logger.Output(2, fmt.Sprintln("DEBUG:   "+lgmsg))
		},
		Trace: func(msg ...interface{}) {
			buf := iorw.NewBuffer()
			buf.Print(msg...)
			lgmsg := buf.String()
			buf = nil
			//log.Println("TRACE: ", lgmsg)
			logger.Output(2, fmt.Sprintln("RACE:   "+lgmsg))
		},
	})
}

var logger *log.Logger = nil

type readwrite struct {
	inbf     *iorw.Buffer
	outw     io.Writer
	lckinout *sync.RWMutex
	flshdr   time.Duration
}

func (rw *readwrite) Write(p []byte) (n int, err error) {
	if pl := len(p); pl > 0 {
		func() {
			rw.lckinout.RLock()
			defer rw.lckinout.RUnlock()
			n, err = rw.inbf.Write(p[:pl])
		}()
	}
	return
}

func (rw *readwrite) ticFlushing() {
	var checksize = func() int64 {
		rw.lckinout.RLock()
		defer rw.lckinout.RUnlock()
		return rw.inbf.Size()
	}
	var bufout = iorw.NewBuffer()
	for {
		time.Sleep(rw.flshdr)
		if checksize() > 0 {
			func() {
				rw.lckinout.Lock()
				defer rw.lckinout.Unlock()
				rd := rw.inbf.Reader()
				io.Copy(bufout, rd)
				rd.Close()
				rw.inbf.Clear()
			}()
			if bufout.Size() > 0 {
				rd := bufout.Reader()
				io.Copy(rw.outw, rd)
				rd.Close()
				bufout.Clear()
			}
		}
	}
}

func init() {
	var rw = &readwrite{inbf: iorw.NewBuffer(), lckinout: &sync.RWMutex{}, flshdr: time.Millisecond * 500, outw: os.Stderr}
	go rw.ticFlushing()
	logger = log.New(rw, "", log.LstdFlags)
}
